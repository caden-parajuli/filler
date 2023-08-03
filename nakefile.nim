import nake
import std/[strutils, os, terminal]

proc `/`[T: openArray[string]](str: string, arr: T): T =
  for i in 0 .. arr.high:
    result[i] = joinPath(str, arr[i])

proc `&`[T: openArray[string]](str: string, arr: T): T =
  for i in 0 .. arr.high:
    result[i] = str & arr[i]
proc `&`[T: openArray[string]](arr: T, str: string): T =
  for i in 0 .. arr.high:
    result[i] = arr[i] & str

const
  srcDir = "src"
  tmpDir = "tmp"
  binDir = "bin"
  nginxDir = "/var/www/rog-ubuntu"
  defaultFlags = "--threads:on --mm:orc"
  releaseFlags = "-d:release"
  debugFlags = ""
  tailwindExe = "./tailwindcss"
  tailwindDebugFlags = ""
  tailwindReleaseFlags = "--minify"
  serverDir = "server"
  serverBin = binDir / "web_socket_server"
  serverBinDir = binDir / serverDir
  
  files = ["web_socket_server"]
  htmlFiles = [""]
  cssFiles = [serverDir / "filler"]
  karaxFiles = [serverDir / "filler"]
  nimJSFiles = [serverDir / "filler"]

proc src(file: string, ext: string = ""): string {.inline.} =
  result = srcDir / (file & ext)
  
proc bin(file: string, ext: string = ""): string {.inline.} =
  result = binDir / (file & ext)

proc tmp(file: string, ext: string = ""): string {.inline.} =
  result = tmpDir / (file & ext)

template genericBuild(debug: bool = true) =
  var
    flags: string
    tailwindFlags: string
  if debug:
    flags = debugFlags
    tailwindFlags = tailwindDebugFlags
  else:
    flags = releaseFlags
    tailwindFlags = tailwindReleaseFlags

  for file in karaxFiles:
    createDir(file.tmp.parentDir)
    createDir(file.bin.parentDir)
    if file != "" and needsRefresh(file.tmp, file.src(".nim")):
      discard direShell(nimExe, "c", defaultFlags, flags, "-o:" & file.tmp, file.src(".nim"))
    if needsRefresh(file.bin(".html"), file.tmp):
      discard direShell(file.tmp, ">", file.bin(".html"))

  for file in htmlFiles:
    createDir(file.bin.parentDir)
    if file != "" and needsRefresh(file.bin(".html"), file.src(".html")):
      copyFile(file.src(".html"), file.bin(".html"))
    
  for file in files:
    createDir(file.bin.parentDir)
    if file != "" and needsRefresh(file.bin, file.src(".nim")):
      discard direShell(nimExe, "c", defaultFlags, flags, "-o:" & file.bin, file.src(".nim"))
  
  for file in cssFiles:
    createDir(file.bin.parentDir)
    if file != "" and needsRefresh(file.bin(".css"), file.src(".css")):
      discard direShell(tailwindExe, "-i", file.src(".css"), "-o", file.bin(".css"), tailwindFlags)

  for file in nimJSFiles:
    createDir(file.bin.parentDir)
    if file != "" and needsRefresh(file.bin(".js"), file.src(".nim")):
      discard direShell(nimExe, "js", "-o:" & file.bin(".js"), file.src(".nim"))

# Build
task "release-build", "Builds everything with release flags":
  genericBuild(debug = false)

task "release", "Builds everything with release flags":
  runTask("releaseBuild")
  
task "debug-build", "Builds everything with debug flags":
  genericBuild(debug = true)
  
task "build", "Runs debug-build":
  runTask("debug-build")

# Clean
task "clean", "Deletes bin and tmp folders":
  removeDir(tmpDir)
  removeDir(binDir)

# Serve
task "debug-serve", "Builds the server in debug mode and copies the files for nginx":
  runTask("debug-build")
  let password = readPasswordFromStdin()
  for file in walkDir(serverBinDir, relative = true):
    discard direShell("echo", password, "| sudo -S", "cp -R", getCurrentDir() / serverBinDir / file.path, nginxDir / file.path)

task "release-serve", "Builds the server in debug mode and copies the files for nginx":
  runTask("release-build")
  let password = readPasswordFromStdin()
  for file in walkDir(serverBinDir, relative = true):
    discard direShell("echo", password, "| sudo -S", "cp -R", getCurrentDir() / serverBinDir / file.path, nginxDir / file.path)

task "serve", "Builds the server and copies the files for nginx":
  runTask("debug-serve")

# Default
task defaultTask, "Default task, runs debug build":
  runTask("debug-build")

