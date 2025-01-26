export type BarcodeStyle = {
    width: string
    height: string
    rotate: string
}

export type ResultColumn = {
    id: number
	test: string
	result: string
	unit: string
	reference_range: string
	abnormal: number
    isNew: boolean
}