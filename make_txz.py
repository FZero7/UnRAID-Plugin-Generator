import tarfile
import hashlib
from pathlib import Path

pkg_dir = Path("smsups-package")
out_file = Path("smsups-package.txz")

# Only the binaries go in the .txz — all text files are embedded in the PLG as INLINE CDATA
BINARY_DIRS = {"usr/local/bin"}

def is_binary(arcname):
    return any(arcname.startswith(d) for d in BINARY_DIRS)

with tarfile.open(out_file, "w:xz") as tar:
    for file_path in sorted(pkg_dir.rglob("*")):
        if file_path.is_file():
            arcname = str(file_path.relative_to(pkg_dir)).replace("\\", "/")
            if not is_binary(arcname):
                continue
            info = tar.gettarinfo(str(file_path), arcname=arcname)
            info.uid = 0
            info.gid = 0
            info.uname = "root"
            info.gname = "root"
            info.mode = 0o755
            with open(file_path, "rb") as f:
                tar.addfile(info, f)

md5 = hashlib.md5(out_file.read_bytes()).hexdigest()
print(f"Created: {out_file} ({out_file.stat().st_size:,} bytes)")
print(f"MD5: {md5}")
