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
  bin = "bin"
  defaultFlags = "--threads:on --mm:orc"
  releaseFlags = "-d:release"
  debugFlags = ""

let
  files = ["server"]

proc src(file: string): string {.inline.} =
  result = src / (file & ".nim")
proc bin(file: string): string {.inline.} =
  result = bin / file

template genericBuild(flags: string) =
  for file in files:
    if needsRefresh*(file.bin, file.src)):
      discard shell(nimExe, "c", flags, file.src, "-o:" & file.bin)
  
task "release-build", "Builds everything with release flags":
  genericBuild(releaseFlags)
task "debug-build", "Builds everything with debug flags":
  genericBuild(debugFlags)
task "build", "Runs debug-build":
  runTask("debug-build")
task defaultTask, "Default task, runs debug build":
  runTask("debug-build")