import tarfile
import hashlib
from pathlib import Path

pkg_dir = Path("smsups-package")
out_file = Path("smsups-package.txz")

def get_mode(arcname):
    if "rc.d" in arcname or "/bin/" in arcname:
        return 0o755
    return 0o644

with tarfile.open(out_file, "w:xz") as tar:
    for file_path in sorted(pkg_dir.rglob("*")):
        if file_path.is_file():
            arcname = str(file_path.relative_to(pkg_dir)).replace("\\", "/")
            info = tar.gettarinfo(str(file_path), arcname=arcname)
            info.uid = 0
            info.gid = 0
            info.uname = "root"
            info.gname = "root"
            info.mode = get_mode(arcname)
            with open(file_path, "rb") as f:
                tar.addfile(info, f)

md5 = hashlib.md5(out_file.read_bytes()).hexdigest()
print(f"Created: {out_file} ({out_file.stat().st_size:,} bytes)")
print(f"MD5: {md5}")
