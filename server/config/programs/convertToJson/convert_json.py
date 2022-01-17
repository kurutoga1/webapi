import os
import sys
import shutil

print(sys.argv)
infile = sys.argv[1]
print(infile)
output_dir = sys.argv[2]
print(output_dir)

print("process start")
# time.sleep(0)
print("process end")

# raise BaseException("python is error")

print("move start")
outfile = os.path.join(output_dir, os.path.basename(infile))
shutil.copy(infile, outfile + ".json")
print("move end")
