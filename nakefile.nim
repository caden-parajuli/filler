import nake
import std/[strutils, os]

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
  src = "src"
  tmp = "tmp"
  bin = "bin"
  defaultFlags = "--threads:on --mm:orc"
  releaseFlags = "-d:release"
  debugFlags = ""
  tailwindExe = "./tailwindcss"
  tailwindDebugFlags = ""
  tailwindReleaseFlags = "--minify"

let
  files: array[string] = ["server/web_socket_server"]
  cssFiles: array[string] = []
  nimJSFiles: array[string] = ["server/index"]
  htmlFiles: array[string] = ["server/index"]
  karaxFiles: array[string] = ["server/index"]

proc src(file: string, ext: string = ""): string {.inline.} =
  result = src / (file & ext)
proc bin(file: string, ext: string = "")): string {.inline.} =
  result = bin / (file & ext)
proc tmp(file: string, ext: string = "")): string {.inline.} =
  result = tmp / (file & ext)

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
    createDir(file.bin.parentDir)
    createDir(file.bin.parentDir)
    if needsRefresh*(file.tmp, file.src(".nim"))):
      discard shell(nimExe, "c", defaultFlags, flags, file.src(".nim"), "-o:" & file.tmp)
      discard shell(file.tmp, ">", file.bin(".html"))

  for file in htmlFiles:
    createDir(file.bin.parentDir)
    if needsRefresh*(file.bin(".html"), file.src(".html")):
      discard copyFile(file.src(".html"), file.bin(".html"))
    
  for file in files:
    createDir(file.bin.parentDir)
    if needsRefresh*(file.bin, file.src(".nim"))):
      discard shell(nimExe, "c", defaultFlags, flags, file.src(".nim"), "-o:" & file.bin)
  
  for file in cssFiles:
    createDir(file.bin.parentDir)
    if needsRefresh*(file.bin, file.src(".css")):
      discard shell(tailwindExe, "-i", file.src(".css"), "-o" & file.bin, tailwindFlags)

  for file in nimJSFiles:
    createDir(file.bin.parentDir)
    if needsRefresh*(file.bin(".js"), file.src(".nim"):
      discard shell(nimExe, "c", file.src(".nim"), "-o:" & file.bin("js"))


task "release-build", "Builds everything with release flags":
  genericBuild(debug = false)

task "release", "Builds everything with release flags":
  runTask("releaseBuild")
  
task "debug-build", "Builds everything with debug flags":
  genericBuild(debug = true)
  
task "build", "Runs debug-build":
  runTask("debug-build")
  
task defaultTask, "Default task, runs debug build":
  runTask("debug-build")
  