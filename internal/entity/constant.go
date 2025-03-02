package entity

const (
	// HeaderXTotalCount is a header for react admin. Expected value is a number with a total of the results
	HeaderXTotalCount = "X-Total-Count"
)

type SpecimenType string

const (
	// 	Abscess
	SpecimenTypeABS SpecimenType = "ABS"
	// Tissue, Acne
	SpecimenTypeACNE SpecimenType = "ACNE"
	// Fluid, Acne
	SpecimenTypeACNFLD SpecimenType = "ACNFLD"
	// Air Sample
	SpecimenTypeAIRS SpecimenType = "AIRS"
	// Allograft
	SpecimenTypeALL SpecimenType = "ALL"
	// Amputation
	SpecimenTypeAMP SpecimenType = "AMP"
	// Catheter Tip, Angio
	SpecimenTypeANGI SpecimenType = "ANGI"
	// Catheter Tip, Arterial
	SpecimenTypeARTC SpecimenType = "ARTC"
	// Serum, Acute
	SpecimenTypeASERU SpecimenType = "ASERU"
	// Aspirate
	SpecimenTypeASP SpecimenType = "ASP"
	// Environment, Attest
	SpecimenTypeATTE SpecimenType = "ATTE"
	// Environmental, Autoclave Ampule
	SpecimenTypeAUTOA SpecimenType = "AUTOA"
	// Environmental, Autoclave Capsule
	SpecimenTypeAUTOC SpecimenType = "AUTOC"
	// Autopsy
	SpecimenTypeAUTP SpecimenType = "AUTP"
	// Blood bag
	SpecimenTypeBBL SpecimenType = "BBL"
	// Cyst, Baker's
	SpecimenTypeBCYST SpecimenType = "BCYST"
	// Bite
	SpecimenTypeBITE SpecimenType = "BITE"
	// Bleb
	SpecimenTypeBLEB SpecimenType = "BLEB"
	// Blister
	SpecimenTypeBLIST SpecimenType = "BLIST"
	// Boil
	SpecimenTypeBOIL SpecimenType = "BOIL"
	// Bone
	SpecimenTypeBON SpecimenType = "BON"
	// Bowel contents
	SpecimenTypeBOWL SpecimenType = "BOWL"
	// Blood product unit
	SpecimenTypeBPU SpecimenType = "BPU"
	// Burn
	SpecimenTypeBRN SpecimenType = "BRN"
	// Brush
	SpecimenTypeBRSH SpecimenType = "BRSH"
	// Breath (use EXHLD)
	SpecimenTypeBRTH SpecimenType = "BRTH"
	// Brushing
	SpecimenTypeBRUS SpecimenType = "BRUS"
	// Bubo
	SpecimenTypeBUB SpecimenType = "BUB"
	// Bulla/Bullae
	SpecimenTypeBULLA SpecimenType = "BULLA"
	// Biopsy
	SpecimenTypeBX SpecimenType = "BX"
	// Calculus (=Stone)
	SpecimenTypeCALC SpecimenType = "CALC"
	// Carbuncle
	SpecimenTypeCARBU SpecimenType = "CARBU"
	// Catheter
	SpecimenTypeCAT SpecimenType = "CAT"
	// Bite, Cat
	SpecimenTypeCBITE SpecimenType = "CBITE"
	// Clippings
	SpecimenTypeCLIPP SpecimenType = "CLIPP"
	// Conjunctiva
	SpecimenTypeCNJT SpecimenType = "CNJT"
	// Colostrum
	SpecimenTypeCOL SpecimenType = "COL"
	// Biospy, Cone
	SpecimenTypeCONE SpecimenType = "CONE"
	// Scratch, Cat
	SpecimenTypeCSCR SpecimenType = "CSCR"
	// Serum, Convalescent
	SpecimenTypeCSERU SpecimenType = "CSERU"
	// Catheter Insertion Site
	SpecimenTypeCSITE SpecimenType = "CSITE"
	// Fluid, Cystostomy Tube
	SpecimenTypeCSMY SpecimenType = "CSMY"
	// Fluid, Cyst
	SpecimenTypeCST SpecimenType = "CST"
	// Blood, Cell Saver
	SpecimenTypeCSVR SpecimenType = "CSVR"
	// Catheter tip
	SpecimenTypeCTP SpecimenType = "CTP"
	// Site, CVP
	SpecimenTypeCVPS SpecimenType = "CVPS"
	// Catheter Tip, CVP
	SpecimenTypeCVPT SpecimenType = "CVPT"
	// Nodule, Cystic
	SpecimenTypeCYN SpecimenType = "CYN"
	// Cyst
	SpecimenTypeCYST SpecimenType = "CYST"
	// Bite, Dog
	SpecimenTypeDBITE SpecimenType = "DBITE"
	// Sputum, Deep Cough
	SpecimenTypeDCS SpecimenType = "DCS"
	// Ulcer, Decubitus
	SpecimenTypeDEC SpecimenType = "DEC"
	// Environmental, Water (Deionized)
	SpecimenTypeDEION SpecimenType = "DEION"
	// Dialysate
	SpecimenTypeDIA SpecimenType = "DIA"
	// Discharge
	SpecimenTypeDISCHG SpecimenType = "DISCHG"
	// Diverticulum
	SpecimenTypeDIV SpecimenType = "DIV"
	// Drain
	SpecimenTypeDRN SpecimenType = "DRN"
	// Drainage, Tube
	SpecimenTypeDRNG SpecimenType = "DRNG"
	// Drainage, Penrose
	SpecimenTypeDRNGP SpecimenType = "DRNGP"
	// Ear wax (cerumen)
	SpecimenTypeEARW SpecimenType = "EARW"
	// Brush, Esophageal
	SpecimenTypeEBRUSH SpecimenType = "EBRUSH"
	// Environmental, Eye Wash
	SpecimenTypeEEYE SpecimenType = "EEYE"
	// Environmental, Effluent
	SpecimenTypeEFF SpecimenType = "EFF"
	// Effusion
	SpecimenTypeEFFUS SpecimenType = "EFFUS"
	// Environmental, Food
	SpecimenTypeEFOD SpecimenType = "EFOD"
	// Environmental, Isolette
	SpecimenTypeEISO SpecimenType = "EISO"
	// Electrode
	SpecimenTypeELT SpecimenType = "ELT"
	// Environmental, Unidentified Substance
	SpecimenTypeENVIR SpecimenType = "ENVIR"
	// Environmental, Other Substance
	SpecimenTypeEOTH SpecimenType = "EOTH"
	// Environmental, Soil
	SpecimenTypeESOI SpecimenType = "ESOI"
	// Environmental, Solution (Sterile)
	SpecimenTypeESOS SpecimenType = "ESOS"
	// Aspirate, Endotrach
	SpecimenTypeETA SpecimenType = "ETA"
	// Catheter Tip, Endotracheal
	SpecimenTypeETTP SpecimenType = "ETTP"
	// Tube, Endotracheal
	SpecimenTypeETTUB SpecimenType = "ETTUB"
	// Environmental, Whirlpool
	SpecimenTypeEWHI SpecimenType = "EWHI"
	// Gas, exhaled (=breath)
	SpecimenTypeEXG SpecimenType = "EXG"
	// Shunt, External
	SpecimenTypeEXS SpecimenType = "EXS"
	// Exudate
	SpecimenTypeEXUDTE SpecimenType = "EXUDTE"
	// Environmental, Water (Well)
	SpecimenTypeFAW SpecimenType = "FAW"
	// Blood, Fetal
	SpecimenTypeFBLOOD SpecimenType = "FBLOOD"
	// Fluid, Abdomen
	SpecimenTypeFGA SpecimenType = "FGA"
	// Fistula
	SpecimenTypeFIST SpecimenType = "FIST"
	// Fluid, Other
	SpecimenTypeFLD SpecimenType = "FLD"
	// Filter
	SpecimenTypeFLT SpecimenType = "FLT"
	// Fluid, Body unsp
	SpecimenTypeFLU SpecimenType = "FLU"
	// Fluid
	SpecimenTypeFLUID SpecimenType = "FLUID"
	// Catheter Tip, Foley
	SpecimenTypeFOLEY SpecimenType = "FOLEY"
	// Fluid, Respiratory
	SpecimenTypeFRS SpecimenType = "FRS"
	// Scalp, Fetal
	SpecimenTypeFSCLP SpecimenType = "FSCLP"
	// Furuncle
	SpecimenTypeFUR SpecimenType = "FUR"
	// Gas
	SpecimenTypeGAS SpecimenType = "GAS"
	// Aspirate, Gastric
	SpecimenTypeGASA SpecimenType = "GASA"
	// Antrum, Gastric
	SpecimenTypeGASAN SpecimenType = "GASAN"
	// Brushing, Gastric
	SpecimenTypeGASBR SpecimenType = "GASBR"
	// Drainage, Gastric
	SpecimenTypeGASD SpecimenType = "GASD"
	// Fluid/contents, Gastric
	SpecimenTypeGAST SpecimenType = "GAST"
	// Genital vaginal
	SpecimenTypeGENV SpecimenType = "GENV"
	// Graft
	SpecimenTypeGRAFT SpecimenType = "GRAFT"
	// Graft Site
	SpecimenTypeGRAFTS SpecimenType = "GRAFTS"
	// Granuloma
	SpecimenTypeGRANU SpecimenType = "GRANU"
	// Catheter, Groshong
	SpecimenTypeGROSH SpecimenType = "GROSH"
	// Solution, Gastrostomy
	SpecimenTypeGSOL SpecimenType = "GSOL"
	// Biopsy, Gastric
	SpecimenTypeGSPEC SpecimenType = "GSPEC"
	// Tube, Gastric
	SpecimenTypeGT SpecimenType = "GT"
	// Drainage Tube, Drainage (Gastrostomy)
	SpecimenTypeGTUBE SpecimenType = "GTUBE"
	// Bite, Human
	SpecimenTypeHBITE SpecimenType = "HBITE"
	// Blood, Autopsy
	SpecimenTypeHBLUD SpecimenType = "HBLUD"
	// Catheter Tip, Hemaquit
	SpecimenTypeHEMAQ SpecimenType = "HEMAQ"
	// Catheter Tip, Hemovac
	SpecimenTypeHEMO SpecimenType = "HEMO"
	// Tissue, Herniated
	SpecimenTypeHERNI SpecimenType = "HERNI"
	// Drain, Hemovac
	SpecimenTypeHEV SpecimenType = "HEV"
	// Catheter, Hickman
	SpecimenTypeHIC SpecimenType = "HIC"
	// Fluid, Hydrocele
	SpecimenTypeHYDC SpecimenType = "HYDC"
	// Bite, Insect
	SpecimenTypeIBITE SpecimenType = "IBITE"
	// Cyst, Inclusion
	SpecimenTypeICYST SpecimenType = "ICYST"
	// Catheter Tip, Indwelling
	SpecimenTypeIDC SpecimenType = "IDC"
	// Gas, Inhaled
	SpecimenTypeIHG SpecimenType = "IHG"
	// Drainage, Ileostomy
	SpecimenTypeILEO SpecimenType = "ILEO"
	// Source of Specimen Is Illegible
	SpecimenTypeILLEG SpecimenType = "ILLEG"
	// Implant
	SpecimenTypeIMP SpecimenType = "IMP"
	// Site, Incision/Surgical
	SpecimenTypeINCI SpecimenType = "INCI"
	// Infiltrate
	SpecimenTypeINFIL SpecimenType = "INFIL"
	// Insect
	SpecimenTypeINS SpecimenType = "INS"
	// Catheter Tip, Introducer
	SpecimenTypeINTRD SpecimenType = "INTRD"
	// Intubation tube
	SpecimenTypeIT SpecimenType = "IT"
	// Intrauterine Device
	SpecimenTypeIUD SpecimenType = "IUD"
	// Catheter Tip, IV
	SpecimenTypeIVCAT SpecimenType = "IVCAT"
	// Fluid, IV
	SpecimenTypeIVFLD SpecimenType = "IVFLD"
	// Tubing Tip, IV
	SpecimenTypeIVTIP SpecimenType = "IVTIP"
	// Drainage, Jejunal
	SpecimenTypeJEJU SpecimenType = "JEJU"
	// Fluid, Joint
	SpecimenTypeJNTFLD SpecimenType = "JNTFLD"
	// Drainage, Jackson Pratt
	SpecimenTypeJP SpecimenType = "JP"
	// Lavage
	SpecimenTypeKELOI SpecimenType = "KELOI"
	// Fluid, Kidney
	SpecimenTypeKIDFLD SpecimenType = "KIDFLD"
	// Lavage, Bronhial
	SpecimenTypeLAVG SpecimenType = "LAVG"
	// Lavage, Gastric
	SpecimenTypeLAVGG SpecimenType = "LAVGG"
	// Lavage, Peritoneal
	SpecimenTypeLAVGP SpecimenType = "LAVGP"
	// Lavage, Pre-Bronch
	SpecimenTypeLAVPG SpecimenType = "LAVPG"
	// Contact Lens
	SpecimenTypeLENS1 SpecimenType = "LENS1"
	// Contact Lens Case
	SpecimenTypeLENS2 SpecimenType = "LENS2"
	// Lesion
	SpecimenTypeLESN SpecimenType = "LESN"
	// Liquid, Unspecified
	SpecimenTypeLIQ SpecimenType = "LIQ"
	// Liquid, Other
	SpecimenTypeLIQO SpecimenType = "LIQO"
	// Fluid, Lumbar Sac
	SpecimenTypeLSAC SpecimenType = "LSAC"
	// Catheter Tip, Makurkour
	SpecimenTypeMAHUR SpecimenType = "MAHUR"
	// Mass
	SpecimenTypeMASS SpecimenType = "MASS"
	// Blood, Menstrual
	SpecimenTypeMBLD SpecimenType = "MBLD"
	// Mucosa
	SpecimenTypeMUCOS SpecimenType = "MUCOS"
	// Mucus
	SpecimenTypeMUCUS SpecimenType = "MUCUS"
	// Drainage, Nasal
	SpecimenTypeNASDR SpecimenType = "NASDR"
	// Needle
	SpecimenTypeNEDL SpecimenType = "NEDL"
	// Site, Nephrostomy
	SpecimenTypeNEPH SpecimenType = "NEPH"
	// Aspirate, Nasogastric
	SpecimenTypeNGASP SpecimenType = "NGASP"
	// Drainage, Nasogastric
	SpecimenTypeNGAST SpecimenType = "NGAST"
	// Site, Naso/Gastric
	SpecimenTypeNGS SpecimenType = "NGS"
	// Nodule(s)
	SpecimenTypeNODUL SpecimenType = "NODUL"
	// Secretion, Nasal
	SpecimenTypeNSECR SpecimenType = "NSECR"
	// Other
	SpecimenTypeORH SpecimenType = "ORH"
	// Lesion, Oral
	SpecimenTypeORL SpecimenType = "ORL"
	// Source, Other
	SpecimenTypeOTH SpecimenType = "OTH"
	// Pacemaker
	SpecimenTypePACEM SpecimenType = "PACEM"
	// Fluid, Pericardial
	SpecimenTypePCFL SpecimenType = "PCFL"
	// Site, Peritoneal Dialysis
	SpecimenTypePDSIT SpecimenType = "PDSIT"
	// Site, Peritoneal Dialysis Tunnel
	SpecimenTypePDTS SpecimenType = "PDTS"
	// Abscess, Pelvic
	SpecimenTypePELVA SpecimenType = "PELVA"
	// Lesion, Penile
	SpecimenTypePENIL SpecimenType = "PENIL"
	// Abscess, Perianal
	SpecimenTypePERIA SpecimenType = "PERIA"
	// Cyst, Pilonidal
	SpecimenTypePILOC SpecimenType = "PILOC"
	// Site, Pin
	SpecimenTypePINS SpecimenType = "PINS"
	// Site, Pacemaker Insetion
	SpecimenTypePIS SpecimenType = "PIS"
	// Plant Material
	SpecimenTypePLAN SpecimenType = "PLAN"
	// Plasma
	SpecimenTypePLAS SpecimenType = "PLAS"
	// Plasma bag
	SpecimenTypePLB SpecimenType = "PLB"
	// Serum, Peak Level
	SpecimenTypePLEVS SpecimenType = "PLEVS"
	// Drainage, Penile
	SpecimenTypePND SpecimenType = "PND"
	// Polyps
	SpecimenTypePOL SpecimenType = "POL"
	// Graft Site, Popliteal
	SpecimenTypePOPGS SpecimenType = "POPGS"
	// Graft, Popliteal
	SpecimenTypePOPLG SpecimenType = "POPLG"
	// Site, Popliteal Vein
	SpecimenTypePOPLV SpecimenType = "POPLV"
	// Catheter, Porta
	SpecimenTypePORTA SpecimenType = "PORTA"
	// Plasma, Platelet poor
	SpecimenTypePPP SpecimenType = "PPP"
	// Prosthetic Device
	SpecimenTypePROST SpecimenType = "PROST"
	// Plasma, Platelet rich
	SpecimenTypePRP SpecimenType = "PRP"
	// Pseudocyst
	SpecimenTypePSC SpecimenType = "PSC"
	// Wound, Puncture
	SpecimenTypePUNCT SpecimenType = "PUNCT"
	// Pus
	SpecimenTypePUS SpecimenType = "PUS"
	// Pustule
	SpecimenTypePUSFR SpecimenType = "PUSFR"
	// Pus
	SpecimenTypePUST SpecimenType = "PUST"
	// Quality Control
	SpecimenTypeQC3 SpecimenType = "QC3"
	// Urine, Random
	SpecimenTypeRANDU SpecimenType = "RANDU"
	// Bite, Reptile
	SpecimenTypeRBITE SpecimenType = "RBITE"
	// Drainage, Rectal
	SpecimenTypeRECT SpecimenType = "RECT"
	// Abscess, Rectal
	SpecimenTypeRECTA SpecimenType = "RECTA"
	// Cyst, Renal
	SpecimenTypeRENALC SpecimenType = "RENALC"
	// Fluid, Renal Cyst
	SpecimenTypeRENC SpecimenType = "RENC"
	// Respiratory
	SpecimenTypeRES SpecimenType = "RES"
	// Saliva
	SpecimenTypeSAL SpecimenType = "SAL"
	// Tissue, Keloid (Scar)
	SpecimenTypeSCAR SpecimenType = "SCAR"
	// Catheter Tip, Subclavian
	SpecimenTypeSCLV SpecimenType = "SCLV"
	// Abscess, Scrotal
	SpecimenTypeSCROA SpecimenType = "SCROA"
	// Secretion(s)
	SpecimenTypeSECRE SpecimenType = "SECRE"
	// Serum
	SpecimenTypeSER SpecimenType = "SER"
	// Site, Shunt
	SpecimenTypeSHU SpecimenType = "SHU"
	// Fluid, Shunt
	SpecimenTypeSHUNF SpecimenType = "SHUNF"
	// Shunt
	SpecimenTypeSHUNT SpecimenType = "SHUNT"
	// Site
	SpecimenTypeSITE SpecimenType = "SITE"
	// Biopsy, Skin
	SpecimenTypeSKBP SpecimenType = "SKBP"
	// Skin
	SpecimenTypeSKN SpecimenType = "SKN"
	// Mass, Sub-Mandibular
	SpecimenTypeSMM SpecimenType = "SMM"
	// Fluid, synovial (Joint fluid)
	SpecimenTypeSNV SpecimenType = "SNV"
	// Spermatozoa
	SpecimenTypeSPRM SpecimenType = "SPRM"
	// Catheter Tip, Suprapubic
	SpecimenTypeSPRP SpecimenType = "SPRP"
	// Cathether Tip, Suprapubic
	SpecimenTypeSPRPB SpecimenType = "SPRPB"
	// Environmental, Spore Strip
	SpecimenTypeSPS SpecimenType = "SPS"
	// Sputum
	SpecimenTypeSPT SpecimenType = "SPT"
	// Sputum - coughed
	SpecimenTypeSPTC SpecimenType = "SPTC"
	// Sputum - tracheal aspirate
	SpecimenTypeSPTT SpecimenType = "SPTT"
	// Sputum, Simulated
	SpecimenTypeSPUT1 SpecimenType = "SPUT1"
	// Sputum, Inducted
	SpecimenTypeSPUTIN SpecimenType = "SPUTIN"
	// Sputum, Spontaneous
	SpecimenTypeSPUTSP SpecimenType = "SPUTSP"
	// Environmental, Sterrad
	SpecimenTypeSTER SpecimenType = "STER"
	// Stool = Fecal
	SpecimenTypeSTL SpecimenType = "STL"
	// Stone, Kidney
	SpecimenTypeSTONE SpecimenType = "STONE"
	// Abscess, Submandibular
	SpecimenTypeSUBMA SpecimenType = "SUBMA"
	// Abscess, Submaxillary
	SpecimenTypeSUBMX SpecimenType = "SUBMX"
	// Drainage, Sump
	SpecimenTypeSUMP SpecimenType = "SUMP"
	// Suprapubic Tap
	SpecimenTypeSUP SpecimenType = "SUP"
	// Suture
	SpecimenTypeSUTUR SpecimenType = "SUTUR"
	// Catheter Tip, Swan Gantz
	SpecimenTypeSWGZ SpecimenType = "SWGZ"
	// Aspirate, Tracheal
	SpecimenTypeTASP SpecimenType = "TASP"
	// Tissue
	SpecimenTypeTISS SpecimenType = "TISS"
	// Tissue ulcer
	SpecimenTypeTISU SpecimenType = "TISU"
	// Cathether Tip, Triple Lumen
	SpecimenTypeTLC SpecimenType = "TLC"
	// Site, Tracheostomy
	SpecimenTypeTRAC SpecimenType = "TRAC"
	// Transudate
	SpecimenTypeTRANS SpecimenType = "TRANS"
	// Serum, Trough
	SpecimenTypeTSERU SpecimenType = "TSERU"
	// Abscess, Testicular
	SpecimenTypeTSTES SpecimenType = "TSTES"
	// Aspirate, Transtracheal
	SpecimenTypeTTRA SpecimenType = "TTRA"
	// Tubes
	SpecimenTypeTUBES SpecimenType = "TUBES"
	// Tumor
	SpecimenTypeTUMOR SpecimenType = "TUMOR"
	// Smear, Tzanck
	SpecimenTypeTZANC SpecimenType = "TZANC"
	// Source, Unidentified
	SpecimenTypeUDENT SpecimenType = "UDENT"
	// Urine
	SpecimenTypeUR SpecimenType = "UR"
	// Urine clean catch
	SpecimenTypeURC SpecimenType = "URC"
	// Urine, Bladder Washings
	SpecimenTypeURINB SpecimenType = "URINB"
	// Urine, Catheterized
	SpecimenTypeURINC SpecimenType = "URINC"
	// Urine, Midstream
	SpecimenTypeURINM SpecimenType = "URINM"
	// Urine, Nephrostomy
	SpecimenTypeURINN SpecimenType = "URINN"
	// Urine, Pedibag
	SpecimenTypeURINP SpecimenType = "URINP"
	// Urine catheter
	SpecimenTypeURT SpecimenType = "URT"
	// Urine, Cystoscopy
	SpecimenTypeUSCOP SpecimenType = "USCOP"
	// Source, Unspecified
	SpecimenTypeUSPEC SpecimenType = "USPEC"
	// Catheter Tip, Vas
	SpecimenTypeVASTIP SpecimenType = "VASTIP"
)

// TODO: Change this to real data
func (s SpecimenType) Code() string {
	switch s {
	case SpecimenTypeSER:
		return "1"
	case SpecimenTypeUR:
		return "2"
	case SpecimenTypeHYDC:
		return "3"
	case SpecimenTypePLAS:
		return "4"
	}

	return "5"
}

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
