 
from rembg import remove
 
import easygui
 
from PIL import Image
   
import sys
 

def BackgroundRemoving(inputPath,ID ):
    input = Image.open(inputPath)
    output = remove(input)
    filename = "urfuckingshitwasdone" + str(ID) + ".png"
    output.save(filename)
    print(filename)



if __name__ == "__main__":
    if len(sys.argv) < 3:  # Check if at least 3 arguments are passed
        print("Usage: python backgroundRemover.py <inputPath> <ID>")
        sys.exit(1)
    
    inputPath = sys.argv[1]  # First argument: input file path
    ID = sys.argv[2]         # Second argument: ID for the filename
    BackgroundRemoving(inputPath, ID)