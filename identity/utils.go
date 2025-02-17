package identity

// Polynomial used for CRC16-XModem checksum in Stellar's StrKey format
const crc16Poly = 0x1021

// crc16Checksum calculates a CRC16-XModem checksum as used in Stellar's StrKey format
func crc16Checksum(data []byte) [2]byte {
	crc := uint16(0)
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ crc16Poly
			} else {
				crc = crc << 1
			}
		}
	}
	return [2]byte{byte(crc >> 8), byte(crc)}
}
