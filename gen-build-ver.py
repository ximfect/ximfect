commit = "UNKNOWN"
with open(".git/refs/heads/master", "rt") as vf:
    commit = vf.read()[:7]
src = ""
src_ver = ""
with open("src/tool/const.go", "rt") as sf:
    src_ver = sf.read()[124:129]

with open("src/tool/const.go", "wt") as cf:
    cf.write(f"""/* generic CLI constants */

package tool

const (
	// Version represents the current version of ximfect
	Version string = "{src_ver}+{commit}"
)
""")
