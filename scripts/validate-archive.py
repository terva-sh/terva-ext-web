#!/usr/bin/env python3
"""Verify one GoReleaser archive and run native protocol/host-driver tests."""
import hashlib
import json
import os
from pathlib import Path
import platform
import shutil
import subprocess
import sys
import tarfile
import tempfile
import zipfile


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    archive_dir, goos, goarch = sys.argv[1:]
    root = Path(__file__).resolve().parents[1]
    mode = "source" if archive_dir == "--source" else "archive"
    report_path = root / f"validation-report-{mode}-{goos}-{goarch}.json"
    report = {
        "commit": subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=root, text=True).strip(),
        "platform": f"{goos}/{goarch}", "native_machine": platform.machine(),
        "go": subprocess.check_output(["go", "version"], text=True).strip(),
        "host_driver": "terva.sh/terva v0.137.0 (floor/current)",
        "mode": mode, "checks": [],
        "runner_image": os.environ.get("ImageOS", "local"),
        "runner_image_version": os.environ.get("ImageVersion", "local"),
        "workflow_run": os.environ.get("GITHUB_RUN_ID", "local"),
    }
    try:
        actual = subprocess.check_output(["go", "env", "GOOS", "GOARCH"], text=True).split()
        assert actual == [goos, goarch], "runtime target mismatch"
        if mode == "source":
            commands = [
                (["go", "vet", "./..."], root),
                (["go", "test", "-race", "-tags", "conformance", "-count=1", "./..."], root),
                (["go", "test", "-race", "-count=1", "./..."], root / "tests/host-contract"),
            ]
            for command, directory in commands:
                subprocess.run(command, cwd=directory, check=True, timeout=300)
                report["checks"].append(" ".join(command) + " in " + str(directory.relative_to(root)))
            report["status"] = "passed"
            return
        suffix = f"_{goos}_{goarch}." + ("zip" if goos == "windows" else "tar.gz")
        matches = list(Path(archive_dir).glob("*" + suffix))
        assert len(matches) == 1, "expected exactly one matching archive"
        archive = matches[0]
        checksums = {}
        for line in (Path(archive_dir) / "checksums.txt").read_text().splitlines():
            checksum, filename = line.split(maxsplit=1)
            checksums[filename.lstrip("*")] = checksum
        assert digest(archive) == checksums[archive.name], "archive checksum mismatch"
        report.update(archive=archive.name, archive_sha256=digest(archive))
        report["checks"].append("archive checksum")
        with tempfile.TemporaryDirectory(prefix="web archive fixture ") as temporary:
            install = Path(temporary) / "installation with spaces"
            install.mkdir()
            if goos == "windows":
                with zipfile.ZipFile(archive) as bundle:
                    bundle.extractall(install)
            else:
                with tarfile.open(archive) as bundle:
                    bundle.extractall(install, filter="data")
            manifest = json.loads((install / "extension.json").read_text())
            assert manifest["name"] == "web"
            assert (install / "skills/web-research/SKILL.md").is_file(), "missing research skill"
            binary = install / ("terva-ext-web.exe" if goos == "windows" else "terva-ext-web")
            report.update(binary_sha256=digest(binary), version=manifest["version"])
            report["files"] = sorted(str(p.relative_to(install)) for p in install.rglob("*") if p.is_file())
            # Supply only the launcher's shell utilities; Go is deliberately absent.
            utilities = Path(temporary) / "launcher utilities"
            utilities.mkdir()
            for name in ["bash", "dirname", "uname", "find"]:
                executable = shutil.which(name)
                assert executable, f"missing launcher utility {name}"
                target = utilities / Path(executable).name
                if os.name == "nt":
                    # Git Bash programs need their adjacent runtime DLLs. Use
                    # their original directory, which must not contain Go.
                    continue
                target.symlink_to(executable)
            env = dict(os.environ)
            if os.name == "nt":
                env["PATH"] = os.pathsep.join(dict.fromkeys(str(Path(shutil.which(n)).parent) for n in ["bash", "dirname", "uname", "find"]))
            else:
                env["PATH"] = str(utilities)
            assert shutil.which("go", path=env["PATH"]) is None, "Go still available during archive launch"
            banner = subprocess.check_output([shutil.which("bash"), str(install / "run.sh"), "--version"], env=env, text=True, timeout=30)
            assert f"terva-ext-web {manifest['version']} " in banner, "binary/manifest version mismatch"
            report["checks"].append("launcher without Go, path with spaces, version and skill")
            test_env = dict(os.environ, WEB_CONFORMANCE_BINARY=str(binary), WEB_CONFORMANCE_INSTALL=str(install))
            subprocess.run(["go", "test", "-tags", "conformance", "-run", "Conformance", "-count=1", "."], cwd=root, env=test_env, check=True, timeout=180)
            report["checks"].append("archive subprocess conformance")
            subprocess.run(["go", "test", "-race", "-count=1", "./..."], cwd=root / "tests/host-contract", env=test_env, check=True, timeout=180)
            report["checks"].append("published host-driver archive launch and permission contract")
        report["status"] = "passed"
    except Exception as error:
        report["status"] = "failed"
        report["error"] = str(error)
        raise
    finally:
        report_path.write_text(json.dumps(report, indent=2) + "\n")


if __name__ == "__main__":
    main()
