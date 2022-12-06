#!/usr/bin/python
# -*- coding: UTF-8

def sand_and_clay_to_ka5_texture(sand, clay):
    "get a rough KA5 soil texture class from given sand and soil content"
    silt = 1.0 - sand - clay
    soil_texture = ""

# fix me - this setup has errors
    if silt < 0.1 and clay < 0.05:
        soil_texture = "SS "
    elif silt < 0.25 and clay < 0.05:
        soil_texture = "SU2"
    elif silt < 0.25 and clay < 0.08:
        soil_texture = "SL2"
    elif silt < 0.40 and clay < 0.08:
        soil_texture = "SU3"
    elif silt < 0.50 and clay < 0.08:
        soil_texture = "SU4"
    elif silt < 0.8 and clay < 0.08:
        soil_texture = "US "
    elif silt >= 0.8 and clay < 0.08:
        soil_texture = "UU "
    elif silt < 0.1 and clay < 0.17:
        soil_texture = "ST2"
    elif silt < 0.4 and clay < 0.12:
        soil_texture = "SL3"
    elif silt < 0.4 and clay < 0.17:
        soil_texture = "SL4"
    elif silt < 0.5 and clay < 0.17:
        soil_texture = "SLU"
    elif silt < 0.65 and clay < 0.17:
        soil_texture = "ULS"
    elif silt >= 0.65 and clay < 0.12:
        soil_texture = "UT2"
    elif silt >= 0.65 and clay < 0.17:
        soil_texture = "UT3"
    elif silt < 0.15 and clay < 0.25:
        soil_texture = "ST3"
    elif silt < 0.30 and clay < 0.25:
        soil_texture = "LS4"
    elif silt < 0.40 and clay < 0.25:
        soil_texture = "LS3"
    elif silt < 0.50 and clay < 0.25:
        soil_texture = "LS2"
    elif silt < 0.65 and clay < 0.30:
        soil_texture = "LU "
    elif silt >= 0.65 and clay < 0.25:
        soil_texture = "UT4"
    elif silt < 0.15 and clay < 0.35:
        soil_texture = "TS4"
    elif silt < 0.30 and clay < 0.45:
        soil_texture = "LTS"
    elif silt < 0.50 and clay < 0.35:
        soil_texture = "LT2"
    elif silt < 0.65 and clay < 0.45:
        soil_texture = "TU3"
    elif silt >= 0.65 and clay >= 0.25:
        soil_texture = "TU4"
    elif silt < 0.15 and clay < 0.45:
        soil_texture = "TS3"
    elif silt < 0.50 and clay < 0.45:
        soil_texture = "LT3"
    elif silt < 0.15 and clay < 0.65:
        soil_texture = "TS2"
    elif silt < 0.30 and clay < 0.65:
        soil_texture = "TL "
    elif silt >= 0.30 and clay < 0.65:
        soil_texture = "TU2"
    elif clay >= 0.65:
        soil_texture = "TT "
    else:
        soil_texture = ""

    return soil_texture

def getPoreVolume(bulkDensity) :
    return 1 - ((bulkDensity/1000) / 2.65)

def getBulkDensityClass(bulkDensity) :
    bulkDensityClass = 1
    bd = bulkDensity / 1000
    if bd < 1.3 :
        bulkDensityClass = 1
    elif bd < 1.5 :
        bulkDensityClass = 2
    elif bd < 1.7 :
        bulkDensityClass = 3
    elif bd < 1.85 :
        bulkDensityClass = 4
    else :
        bulkDensityClass = 5
    return bulkDensityClass

# PTF nach Toth 2015
#FK:  Let W(LT)    = 0.2449 - 0.1887 * (1/(CGEHALT(1)+1)) + 0.004527 * Ton(1) + 0.001535 * SLUF(1) + 0.001442 * SLUF(1) * (1/(CGEHALT(1)+1)) - 0.0000511 * SLUF(1) * Ton(1) + 0.0008676 * Ton(1) * (1/(CGEHALT(1)+1))
def calcFK(cgehalt, ton, sluf ) :
    return  0.2449 - 0.1887 * (1/(cgehalt+1)) + 0.004527 * ton + 0.001535 * sluf + 0.001442 * sluf * (1/(cgehalt+1)) - 0.0000511 * sluf * ton + 0.0008676 * ton * (1/(cgehalt+1))
# PWP: Let WMIN(LT) = 0.09878 + 0.002127* Ton(1) - 0.0008366 *SLUF(1) - 0.0767 *(1/(CGEHALT(1)+1)) + 0.00003853 * SLUF(1) * Ton(1) + 0.00233 * SLUF(1) * (1/(CGEHALT(1)+1)) + 0.0009498 * SLUF(1) * (1/(CGEHALT(1)+1))
def calcPWP(cgehalt, ton, sluf) :
    val = 0.09878 + 0.002127 * ton - 0.0008366 * sluf - 0.0767 * (1/(cgehalt+1)) + 0.00003853 * sluf * ton + 0.00233 * sluf * (1/(cgehalt+1)) + 0.0009498 * sluf * (1/(cgehalt+1))
    return val

if __name__ == "__main__":
    bulkDensity = 1.5
    print("BDclass:", getBulkDensityClass(bulkDensity))
   # 47	20	33
    sand = 47.0
    silt = 33.0
    clay = 20.0
#SID Corg Te  lb B St C/N C/S Hy Rd NuHo FC WP PS S% SI% C% lamda DraiT  Drai% GW LBG
#009 1.05 LS3 03 2 00 10      00 13 02   29 15 46 47 33 20 00  20   00 99 01

    print("Tex:", sand_and_clay_to_ka5_texture(sand/100, clay/100))
    cgehalt = 1.05 
    print("FC: ", calcFK(cgehalt, clay, silt)*100)
    print("PWP: ", calcPWP(cgehalt, clay, silt)*100)
