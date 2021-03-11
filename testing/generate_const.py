def main(argv):
    with open("version") as f:
        ver = f.read()
    with open("build") as f:
        build = int(f.read()) + 1
    out = (
        "package tool\n\n"
        "const (\n"
        "\t// Version is the release number\n"
        f"\tVersion = \"{ver}\"\n\n"
        "\t// Build is the build number\n"
        f"\tBuild = {build}\n"
        ")")
    with open("../src/tool/const.go", "w") as f:
        f.write(out)
    with open("build", "w") as f:
        f.write(str(build))
    
if __name__ == "__main__":
    from sys import argv
    main(argv)
