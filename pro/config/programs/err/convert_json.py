import os
import sys
import shutil
import time

print(sys.argv)
infile = sys.argv[1]
print(infile)
output_dir = sys.argv[2]
print(output_dir)
sleep_time = int(sys.argv[3])

print("process start")
time.sleep(sleep_time)
raise BaseException("error occur!!")
print("process end")

# raise BaseException("python is error")

print("move start")
outfile = os.path.join(output_dir, os.path.basename(infile)) + ".zip"
shutil.copy(infile, outfile)
print("move end")
