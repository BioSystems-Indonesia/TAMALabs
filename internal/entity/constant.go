package entity

const (
	// HeaderXTotalCount is a header for react admin. Expected value is a number with a total of the results
	HeaderXTotalCount = "X-Total-Count"
)

type SpecimenType string

const (
	// 	Abscess
	SpecimenTypeABS = "ABS"
	// Tissue, Acne
	SpecimenTypeACNE = "ACNE"
	// Fluid, Acne
	SpecimenTypeACNFLD = "ACNFLD"
	// Air Sample
	SpecimenTypeAIRS = "AIRS"
	// Allograft
	SpecimenTypeALL = "ALL"
	// Amputation
	SpecimenTypeAMP = "AMP"
	// Catheter Tip, Angio
	SpecimenTypeANGI = "ANGI"
	// Catheter Tip, Arterial
	SpecimenTypeARTC = "ARTC"
	// Serum, Acute
	SpecimenTypeASERU = "ASERU"
	// Aspirate
	SpecimenTypeASP = "ASP"
	// Environment, Attest
	SpecimenTypeATTE = "ATTE"
	// Environmental, Autoclave Ampule
	SpecimenTypeAUTOA = "AUTOA"
	// Environmental, Autoclave Capsule
	SpecimenTypeAUTOC = "AUTOC"
	// Autopsy
	SpecimenTypeAUTP = "AUTP"
	// Blood bag
	SpecimenTypeBBL = "BBL"
	// Cyst, Baker's
	SpecimenTypeBCYST = "BCYST"
	// Bite
	SpecimenTypeBITE = "BITE"
	// Bleb
	SpecimenTypeBLEB = "BLEB"
	// Blister
	SpecimenTypeBLIST = "BLIST"
	// Boil
	SpecimenTypeBOIL = "BOIL"
	// Bone
	SpecimenTypeBON = "BON"
	// Bowel contents
	SpecimenTypeBOWL = "BOWL"
	// Blood product unit
	SpecimenTypeBPU = "BPU"
	// Burn
	SpecimenTypeBRN = "BRN"
	// Brush
	SpecimenTypeBRSH = "BRSH"
	// Breath (use EXHLD)
	SpecimenTypeBRTH = "BRTH"
	// Brushing
	SpecimenTypeBRUS = "BRUS"
	// Bubo
	SpecimenTypeBUB = "BUB"
	// Bulla/Bullae
	SpecimenTypeBULLA = "BULLA"
	// Biopsy
	SpecimenTypeBX = "BX"
	// Calculus (=Stone)
	SpecimenTypeCALC = "CALC"
	// Carbuncle
	SpecimenTypeCARBU = "CARBU"
	// Catheter
	SpecimenTypeCAT = "CAT"
	// Bite, Cat
	SpecimenTypeCBITE = "CBITE"
	// Clippings
	SpecimenTypeCLIPP = "CLIPP"
	// Conjunctiva
	SpecimenTypeCNJT = "CNJT"
	// Colostrum
	SpecimenTypeCOL = "COL"
	// Biospy, Cone
	SpecimenTypeCONE = "CONE"
	// Scratch, Cat
	SpecimenTypeCSCR = "CSCR"
	// Serum, Convalescent
	SpecimenTypeCSERU = "CSERU"
	// Catheter Insertion Site
	SpecimenTypeCSITE = "CSITE"
	// Fluid, Cystostomy Tube
	SpecimenTypeCSMY = "CSMY"
	// Fluid, Cyst
	SpecimenTypeCST = "CST"
	// Blood, Cell Saver
	SpecimenTypeCSVR = "CSVR"
	// Catheter tip
	SpecimenTypeCTP = "CTP"
	// Site, CVP
	SpecimenTypeCVPS = "CVPS"
	// Catheter Tip, CVP
	SpecimenTypeCVPT = "CVPT"
	// Nodule, Cystic
	SpecimenTypeCYN = "CYN"
	// Cyst
	SpecimenTypeCYST = "CYST"
	// Bite, Dog
	SpecimenTypeDBITE = "DBITE"
	// Sputum, Deep Cough
	SpecimenTypeDCS = "DCS"
	// Ulcer, Decubitus
	SpecimenTypeDEC = "DEC"
	// Environmental, Water (Deionized)
	SpecimenTypeDEION = "DEION"
	// Dialysate
	SpecimenTypeDIA = "DIA"
	// Discharge
	SpecimenTypeDISCHG = "DISCHG"
	// Diverticulum
	SpecimenTypeDIV = "DIV"
	// Drain
	SpecimenTypeDRN = "DRN"
	// Drainage, Tube
	SpecimenTypeDRNG = "DRNG"
	// Drainage, Penrose
	SpecimenTypeDRNGP = "DRNGP"
	// Ear wax (cerumen)
	SpecimenTypeEARW = "EARW"
	// Brush, Esophageal
	SpecimenTypeEBRUSH = "EBRUSH"
	// Environmental, Eye Wash
	SpecimenTypeEEYE = "EEYE"
	// Environmental, Effluent
	SpecimenTypeEFF = "EFF"
	// Effusion
	SpecimenTypeEFFUS = "EFFUS"
	// Environmental, Food
	SpecimenTypeEFOD = "EFOD"
	// Environmental, Isolette
	SpecimenTypeEISO = "EISO"
	// Electrode
	SpecimenTypeELT = "ELT"
	// Environmental, Unidentified Substance
	SpecimenTypeENVIR = "ENVIR"
	// Environmental, Other Substance
	SpecimenTypeEOTH = "EOTH"
	// Environmental, Soil
	SpecimenTypeESOI = "ESOI"
	// Environmental, Solution (Sterile)
	SpecimenTypeESOS = "ESOS"
	// Aspirate, Endotrach
	SpecimenTypeETA = "ETA"
	// Catheter Tip, Endotracheal
	SpecimenTypeETTP = "ETTP"
	// Tube, Endotracheal
	SpecimenTypeETTUB = "ETTUB"
	// Environmental, Whirlpool
	SpecimenTypeEWHI = "EWHI"
	// Gas, exhaled (=breath)
	SpecimenTypeEXG = "EXG"
	// Shunt, External
	SpecimenTypeEXS = "EXS"
	// Exudate
	SpecimenTypeEXUDTE = "EXUDTE"
	// Environmental, Water (Well)
	SpecimenTypeFAW = "FAW"
	// Blood, Fetal
	SpecimenTypeFBLOOD = "FBLOOD"
	// Fluid, Abdomen
	SpecimenTypeFGA = "FGA"
	// Fistula
	SpecimenTypeFIST = "FIST"
	// Fluid, Other
	SpecimenTypeFLD = "FLD"
	// Filter
	SpecimenTypeFLT = "FLT"
	// Fluid, Body unsp
	SpecimenTypeFLU = "FLU"
	// Fluid
	SpecimenTypeFLUID = "FLUID"
	// Catheter Tip, Foley
	SpecimenTypeFOLEY = "FOLEY"
	// Fluid, Respiratory
	SpecimenTypeFRS = "FRS"
	// Scalp, Fetal
	SpecimenTypeFSCLP = "FSCLP"
	// Furuncle
	SpecimenTypeFUR = "FUR"
	// Gas
	SpecimenTypeGAS = "GAS"
	// Aspirate, Gastric
	SpecimenTypeGASA = "GASA"
	// Antrum, Gastric
	SpecimenTypeGASAN = "GASAN"
	// Brushing, Gastric
	SpecimenTypeGASBR = "GASBR"
	// Drainage, Gastric
	SpecimenTypeGASD = "GASD"
	// Fluid/contents, Gastric
	SpecimenTypeGAST = "GAST"
	// Genital vaginal
	SpecimenTypeGENV = "GENV"
	// Graft
	SpecimenTypeGRAFT = "GRAFT"
	// Graft Site
	SpecimenTypeGRAFTS = "GRAFTS"
	// Granuloma
	SpecimenTypeGRANU = "GRANU"
	// Catheter, Groshong
	SpecimenTypeGROSH = "GROSH"
	// Solution, Gastrostomy
	SpecimenTypeGSOL = "GSOL"
	// Biopsy, Gastric
	SpecimenTypeGSPEC = "GSPEC"
	// Tube, Gastric
	SpecimenTypeGT = "GT"
	// Drainage Tube, Drainage (Gastrostomy)
	SpecimenTypeGTUBE = "GTUBE"
	// Bite, Human
	SpecimenTypeHBITE = "HBITE"
	// Blood, Autopsy
	SpecimenTypeHBLUD = "HBLUD"
	// Catheter Tip, Hemaquit
	SpecimenTypeHEMAQ = "HEMAQ"
	// Catheter Tip, Hemovac
	SpecimenTypeHEMO = "HEMO"
	// Tissue, Herniated
	SpecimenTypeHERNI = "HERNI"
	// Drain, Hemovac
	SpecimenTypeHEV = "HEV"
	// Catheter, Hickman
	SpecimenTypeHIC = "HIC"
	// Fluid, Hydrocele
	SpecimenTypeHYDC = "HYDC"
	// Bite, Insect
	SpecimenTypeIBITE = "IBITE"
	// Cyst, Inclusion
	SpecimenTypeICYST = "ICYST"
	// Catheter Tip, Indwelling
	SpecimenTypeIDC = "IDC"
	// Gas, Inhaled
	SpecimenTypeIHG = "IHG"
	// Drainage, Ileostomy
	SpecimenTypeILEO = "ILEO"
	// Source of Specimen Is Illegible
	SpecimenTypeILLEG = "ILLEG"
	// Implant
	SpecimenTypeIMP = "IMP"
	// Site, Incision/Surgical
	SpecimenTypeINCI = "INCI"
	// Infiltrate
	SpecimenTypeINFIL = "INFIL"
	// Insect
	SpecimenTypeINS = "INS"
	// Catheter Tip, Introducer
	SpecimenTypeINTRD = "INTRD"
	// Intubation tube
	SpecimenTypeIT = "IT"
	// Intrauterine Device
	SpecimenTypeIUD = "IUD"
	// Catheter Tip, IV
	SpecimenTypeIVCAT = "IVCAT"
	// Fluid, IV
	SpecimenTypeIVFLD = "IVFLD"
	// Tubing Tip, IV
	SpecimenTypeIVTIP = "IVTIP"
	// Drainage, Jejunal
	SpecimenTypeJEJU = "JEJU"
	// Fluid, Joint
	SpecimenTypeJNTFLD = "JNTFLD"
	// Drainage, Jackson Pratt
	SpecimenTypeJP = "JP"
	// Lavage
	SpecimenTypeKELOI = "KELOI"
	// Fluid, Kidney
	SpecimenTypeKIDFLD = "KIDFLD"
	// Lavage, Bronhial
	SpecimenTypeLAVG = "LAVG"
	// Lavage, Gastric
	SpecimenTypeLAVGG = "LAVGG"
	// Lavage, Peritoneal
	SpecimenTypeLAVGP = "LAVGP"
	// Lavage, Pre-Bronch
	SpecimenTypeLAVPG = "LAVPG"
	// Contact Lens
	SpecimenTypeLENS1 = "LENS1"
	// Contact Lens Case
	SpecimenTypeLENS2 = "LENS2"
	// Lesion
	SpecimenTypeLESN = "LESN"
	// Liquid, Unspecified
	SpecimenTypeLIQ = "LIQ"
	// Liquid, Other
	SpecimenTypeLIQO = "LIQO"
	// Fluid, Lumbar Sac
	SpecimenTypeLSAC = "LSAC"
	// Catheter Tip, Makurkour
	SpecimenTypeMAHUR = "MAHUR"
	// Mass
	SpecimenTypeMASS = "MASS"
	// Blood, Menstrual
	SpecimenTypeMBLD = "MBLD"
	// Mucosa
	SpecimenTypeMUCOS = "MUCOS"
	// Mucus
	SpecimenTypeMUCUS = "MUCUS"
	// Drainage, Nasal
	SpecimenTypeNASDR = "NASDR"
	// Needle
	SpecimenTypeNEDL = "NEDL"
	// Site, Nephrostomy
	SpecimenTypeNEPH = "NEPH"
	// Aspirate, Nasogastric
	SpecimenTypeNGASP = "NGASP"
	// Drainage, Nasogastric
	SpecimenTypeNGAST = "NGAST"
	// Site, Naso/Gastric
	SpecimenTypeNGS = "NGS"
	// Nodule(s)
	SpecimenTypeNODUL = "NODUL"
	// Secretion, Nasal
	SpecimenTypeNSECR = "NSECR"
	// Other
	SpecimenTypeORH = "ORH"
	// Lesion, Oral
	SpecimenTypeORL = "ORL"
	// Source, Other
	SpecimenTypeOTH = "OTH"
	// Pacemaker
	SpecimenTypePACEM = "PACEM"
	// Fluid, Pericardial
	SpecimenTypePCFL = "PCFL"
	// Site, Peritoneal Dialysis
	SpecimenTypePDSIT = "PDSIT"
	// Site, Peritoneal Dialysis Tunnel
	SpecimenTypePDTS = "PDTS"
	// Abscess, Pelvic
	SpecimenTypePELVA = "PELVA"
	// Lesion, Penile
	SpecimenTypePENIL = "PENIL"
	// Abscess, Perianal
	SpecimenTypePERIA = "PERIA"
	// Cyst, Pilonidal
	SpecimenTypePILOC = "PILOC"
	// Site, Pin
	SpecimenTypePINS = "PINS"
	// Site, Pacemaker Insetion
	SpecimenTypePIS = "PIS"
	// Plant Material
	SpecimenTypePLAN = "PLAN"
	// Plasma
	SpecimenTypePLAS = "PLAS"
	// Plasma bag
	SpecimenTypePLB = "PLB"
	// Serum, Peak Level
	SpecimenTypePLEVS = "PLEVS"
	// Drainage, Penile
	SpecimenTypePND = "PND"
	// Polyps
	SpecimenTypePOL = "POL"
	// Graft Site, Popliteal
	SpecimenTypePOPGS = "POPGS"
	// Graft, Popliteal
	SpecimenTypePOPLG = "POPLG"
	// Site, Popliteal Vein
	SpecimenTypePOPLV = "POPLV"
	// Catheter, Porta
	SpecimenTypePORTA = "PORTA"
	// Plasma, Platelet poor
	SpecimenTypePPP = "PPP"
	// Prosthetic Device
	SpecimenTypePROST = "PROST"
	// Plasma, Platelet rich
	SpecimenTypePRP = "PRP"
	// Pseudocyst
	SpecimenTypePSC = "PSC"
	// Wound, Puncture
	SpecimenTypePUNCT = "PUNCT"
	// Pus
	SpecimenTypePUS = "PUS"
	// Pustule
	SpecimenTypePUSFR = "PUSFR"
	// Pus
	SpecimenTypePUST = "PUST"
	// Quality Control
	SpecimenTypeQC3 = "QC3"
	// Urine, Random
	SpecimenTypeRANDU = "RANDU"
	// Bite, Reptile
	SpecimenTypeRBITE = "RBITE"
	// Drainage, Rectal
	SpecimenTypeRECT = "RECT"
	// Abscess, Rectal
	SpecimenTypeRECTA = "RECTA"
	// Cyst, Renal
	SpecimenTypeRENALC = "RENALC"
	// Fluid, Renal Cyst
	SpecimenTypeRENC = "RENC"
	// Respiratory
	SpecimenTypeRES = "RES"
	// Saliva
	SpecimenTypeSAL = "SAL"
	// Tissue, Keloid (Scar)
	SpecimenTypeSCAR = "SCAR"
	// Catheter Tip, Subclavian
	SpecimenTypeSCLV = "SCLV"
	// Abscess, Scrotal
	SpecimenTypeSCROA = "SCROA"
	// Secretion(s)
	SpecimenTypeSECRE = "SECRE"
	// Serum
	SpecimenTypeSER = "SER"
	// Site, Shunt
	SpecimenTypeSHU = "SHU"
	// Fluid, Shunt
	SpecimenTypeSHUNF = "SHUNF"
	// Shunt
	SpecimenTypeSHUNT = "SHUNT"
	// Site
	SpecimenTypeSITE = "SITE"
	// Biopsy, Skin
	SpecimenTypeSKBP = "SKBP"
	// Skin
	SpecimenTypeSKN = "SKN"
	// Mass, Sub-Mandibular
	SpecimenTypeSMM = "SMM"
	// Fluid, synovial (Joint fluid)
	SpecimenTypeSNV = "SNV"
	// Spermatozoa
	SpecimenTypeSPRM = "SPRM"
	// Catheter Tip, Suprapubic
	SpecimenTypeSPRP = "SPRP"
	// Cathether Tip, Suprapubic
	SpecimenTypeSPRPB = "SPRPB"
	// Environmental, Spore Strip
	SpecimenTypeSPS = "SPS"
	// Sputum
	SpecimenTypeSPT = "SPT"
	// Sputum - coughed
	SpecimenTypeSPTC = "SPTC"
	// Sputum - tracheal aspirate
	SpecimenTypeSPTT = "SPTT"
	// Sputum, Simulated
	SpecimenTypeSPUT1 = "SPUT1"
	// Sputum, Inducted
	SpecimenTypeSPUTIN = "SPUTIN"
	// Sputum, Spontaneous
	SpecimenTypeSPUTSP = "SPUTSP"
	// Environmental, Sterrad
	SpecimenTypeSTER = "STER"
	// Stool = Fecal
	SpecimenTypeSTL = "STL"
	// Stone, Kidney
	SpecimenTypeSTONE = "STONE"
	// Abscess, Submandibular
	SpecimenTypeSUBMA = "SUBMA"
	// Abscess, Submaxillary
	SpecimenTypeSUBMX = "SUBMX"
	// Drainage, Sump
	SpecimenTypeSUMP = "SUMP"
	// Suprapubic Tap
	SpecimenTypeSUP = "SUP"
	// Suture
	SpecimenTypeSUTUR = "SUTUR"
	// Catheter Tip, Swan Gantz
	SpecimenTypeSWGZ = "SWGZ"
	// Aspirate, Tracheal
	SpecimenTypeTASP = "TASP"
	// Tissue
	SpecimenTypeTISS = "TISS"
	// Tissue ulcer
	SpecimenTypeTISU = "TISU"
	// Cathether Tip, Triple Lumen
	SpecimenTypeTLC = "TLC"
	// Site, Tracheostomy
	SpecimenTypeTRAC = "TRAC"
	// Transudate
	SpecimenTypeTRANS = "TRANS"
	// Serum, Trough
	SpecimenTypeTSERU = "TSERU"
	// Abscess, Testicular
	SpecimenTypeTSTES = "TSTES"
	// Aspirate, Transtracheal
	SpecimenTypeTTRA = "TTRA"
	// Tubes
	SpecimenTypeTUBES = "TUBES"
	// Tumor
	SpecimenTypeTUMOR = "TUMOR"
	// Smear, Tzanck
	SpecimenTypeTZANC = "TZANC"
	// Source, Unidentified
	SpecimenTypeUDENT = "UDENT"
	// Urine
	SpecimenTypeUR = "UR"
	// Urine clean catch
	SpecimenTypeURC = "URC"
	// Urine, Bladder Washings
	SpecimenTypeURINB = "URINB"
	// Urine, Catheterized
	SpecimenTypeURINC = "URINC"
	// Urine, Midstream
	SpecimenTypeURINM = "URINM"
	// Urine, Nephrostomy
	SpecimenTypeURINN = "URINN"
	// Urine, Pedibag
	SpecimenTypeURINP = "URINP"
	// Urine catheter
	SpecimenTypeURT = "URT"
	// Urine, Cystoscopy
	SpecimenTypeUSCOP = "USCOP"
	// Source, Unspecified
	SpecimenTypeUSPEC = "USPEC"
	// Catheter Tip, Vas
	SpecimenTypeVASTIP = "VASTIP"
)

type Priority string

const (
	// As soon as possible (a priority lower than stat)
	PriorityA Priority = "A"
	// Preoperative (to be done prior to surgery)
	PriorityP Priority = "P"
	// Routine
	PriorityR Priority = "R"
	// Stat (do immediately)
	PriorityS Priority = "S"
	// Timing critical (do as near as possible to requested time)
	PriorityT Priority = "T"
)
