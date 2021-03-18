import subprocess
import os
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
		print("-> Unsuporrted platform!")
		return 2
	out_path = os.path.abspath(out_exec)
	
	print(": Build")
	if os.path.exists(out_path):
		print(": Deleting old executable")
		os.remove(out_path)

	b = subprocess.run(
		join(["go", "build", "-o", out_path]), 
		shell = True, 
		cwd = os.path.abspath("../src"))

	if b.returncode != 0:
		print(b.args)
		print(f"-> Return code: {b.returncode}")
		return 1

	print(f":: Output: {out_path}")

	return 0
	
if __name__ == "__main__":
	from sys import argv
	exit(main(argv) or 0)
