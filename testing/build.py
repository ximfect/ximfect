import os
import subprocess
from shlex import join
from sys import platform


def main(argv) -> int:
	print(": Increment build number")
	p = "../src/tool/build"
	with open(p) as f:
		build = int(f.read()) + 1
	with open(p, "wt") as f:
		f.write(str(build))

	if platform.startswith("linux"):
		out_exec = "ximfect"
	elif platform.startswith("win32") or platform.startswith("cygwin"):
		out_exec = "ximfect.exe"
	else:
		print("!! Unsupported platform!")
		return 2
	out_path = os.path.abspath(out_exec)
	
	print(f": Output: {out_path}")

	if os.path.exists(out_path):
		print(": Deleting old executable")
		os.remove(out_path)

	print(": Build")
	cmd = join(["go", "build", "-o", out_path])
	if platform.startswith("win32"):
		cmd = cmd.replace("'", '"')
		
	b = subprocess.run(
		cmd, 
		shell = True, 
		cwd = os.path.abspath("../src"))

	if b.returncode != 0:
		print(f"!! Return code: {b.returncode}")
		return 1
	
	return 0
	
if __name__ == "__main__":
	from sys import argv
	exit(main(argv) or 0)
