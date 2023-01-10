
#!/usr/bin/python
# -*- coding: UTF-8

latLonFiles = ["missingregions.csv", "gridcells_altitude_ZALF-DK94-DK59.csv"]
outLatLonFile = "gridcells_latlon.csv"

def writeLookupFile() :
    
    lookupGrid = dict()
    with open(outLatLonFile, mode="wt", newline="") as outlookupfile :
        outlookupfile.writelines("GRID_NO,LATITUDE,LONGITUDE,ALTITUDE\n")
        for file in latLonFiles :
            with open(file) as sourcefile:
                firstLine = True
                header = dict()
                for line in sourcefile:
                    if firstLine :
                        firstLine = False
                        header = ReadHeader(line)
                        continue

                    # read relevant content from line 
                    tokens = line.split(",")
                    gridIdx = tokens[header["grid_no"]] 
                    lati = tokens[header["lat"]].strip()
                    longi = tokens[header["lon"]].strip()
                    alti = tokens[header["alti"]].strip()

                    if gridIdx in lookupGrid :
                        continue
                    lookupGrid[gridIdx] = True
                    outline = [
                        gridIdx,
                        lati,
                        longi,
                        alti,
                    ]
                                
                    outlookupfile.writelines(",".join(outline) + "\n")
          
def ReadHeader(line) : 
    #read header
    #"GRID_NO","LATITUDE","LONGITUDE","ALTITUDE","DAY","TEMPERATURE_MAX","TEMPERATURE_MIN","TEMPERATURE_AVG","WINDSPEED","VAPOURPRESSURE","PRECIPITATION","RADIATION"
    #GRID_NO,LATITUDE,LONGITUDE,ALTITUDE
    tokens = line.split(",")
    outDic = dict()
    i = -1
    for token in tokens :
        token = token.strip('\"')
        token = token.strip()
        i = i+1
        if token == "LATITUDE":
            outDic["lat"] = i
        if token == "LONGITUDE":
            outDic["lon"] = i
        if token == "GRID_NO" : 
            outDic["grid_no"] = i
        if token == "ALTITUDE" : 
            outDic["alti"] = i

    return outDic

if __name__ == "__main__":
    writeLookupFile()