package otr4

import (
	"crypto/rand"

	"testing"

	"github.com/twstrike/ed448"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type OTR4Suite struct{}

var _ = Suite(&OTR4Suite{})

func (s *OTR4Suite) Test_Concat(c *C) {
	empty := []byte{}
	bytes := []byte{
		0x04, 0x2a, 0xf3, 0xcc, 0x69, 0xbb, 0xa1, 0x50,
	}

	exp := []byte{
		0x04, 0x2a, 0xf3, 0xcc, 0x69, 0xbb, 0xa1, 0x50,
		0x71, 0x7b, 0x24, 0xd5, 0xd4, 0x98, 0x0c, 0xfe,
		0xce, 0x60, 0xe7, 0x97, 0x84, 0xf4, 0x1c, 0x72,
		0x01, 0x07, 0xb8, 0x24, 0xa8, 0x43, 0x0e, 0x81,
		0x25, 0xca, 0xb4, 0xa0, 0xda, 0xf5, 0xfa, 0xf6,
		0x0c, 0x90, 0x99, 0x7f, 0x1e, 0xed, 0x83, 0xde,
		0xbe, 0xe7, 0xef, 0x8e, 0xea, 0xeb, 0xc8, 0x5d,
		0x67, 0x5b, 0x3b, 0x04, 0x55, 0x0a, 0x36, 0x2f,
		0x06, 0xea, 0x48, 0xc4, 0x23, 0x28, 0xe1, 0x99,
		0x08, 0xa5, 0x88, 0x8f, 0xad, 0x7f, 0x39, 0xdf,
		0x56, 0xa3, 0xaa, 0x4d, 0x59, 0x66, 0xec, 0xd5,
		0x6c, 0x38, 0x02, 0x8c, 0x80, 0x96, 0xd2, 0xd4,
		0x54, 0x24, 0x76, 0x70, 0xda, 0x99, 0xc5, 0xd6,
		0x81, 0x40, 0x49, 0xcd, 0x76, 0xb1, 0x05, 0xc4,
		0xa8, 0x42, 0x17, 0x09, 0x51, 0xc2, 0xa9, 0x2e,
	}

	c.Assert(func() { concat() }, Panics, "programmer error: missing concat arguments")
	c.Assert(func() { concat(bytes) }, Panics, "programmer error: missing concat arguments")
	c.Assert(func() { concat("not a valid input", bytes) }, Panics, "programmer error: invalid input")
	c.Assert(concat(empty, bytes, testSec, testPubA), DeepEquals, exp)
}

func (s *OTR4Suite) Test_Auth(c *C) {
	message := []byte("our message")
	out, err := auth(fixedRand(randAuthData), testPubA, testPubB, testPubC, testSec, message)

	c.Assert(out, DeepEquals, testSigma)
	c.Assert(err, IsNil)

	r := make([]byte, 270)
	out, err = auth(fixedRand(r), testPubA, testPubB, testPubC, testSec, message)

	c.Assert(err, ErrorMatches, ".*cannot source enough entropy")
	c.Assert(out, IsNil)

	r = make([]byte, 56)
	out, err = auth(fixedRand(r), testPubA, testPubB, testPubC, testSec, message)

	c.Assert(err, ErrorMatches, ".*cannot source enough entropy")
	c.Assert(out, IsNil)
}

func (s *OTR4Suite) Test_Verify(c *C) {
	message := []byte("our message")

	b := verify(testPubA, testPubB, testPubC, testSigma, message)

	c.Assert(b, Equals, true)
}

func (s *OTR4Suite) Test_VerifyAndAuth(c *C) {
	message := []byte("hello, I am a message")
	sigma, _ := auth(rand.Reader, testPubA, testPubB, testPubC, testSec, message)
	ver := verify(testPubA, testPubB, testPubC, sigma, message)
	c.Assert(ver, Equals, true)

	fakeMessage := []byte("fake message")
	ver = verify(testPubA, testPubB, testPubC, sigma, fakeMessage)
	c.Assert(ver, Equals, false)

	ver = verify(testPubB, testPubB, testPubC, sigma, message)
	c.Assert(ver, Equals, false)

	ver = verify(testPubA, testPubA, testPubC, sigma, message)
	c.Assert(ver, Equals, false)

	ver = verify(testPubA, testPubB, testPubB, sigma, message)
	c.Assert(ver, Equals, false)

	ver = verify(testPubA, testPubB, testPubC, testSigma, message)
	c.Assert(ver, Equals, false)
}

func (s *OTR4Suite) Test_DREnc(c *C) {
	//XXX: move this data into a file
	pub1 := &cramerShoupPublicKey{
		ed448.NewPoint(
			[16]uint32{
				0x04928f75, 0x086d49d0, 0x01204b29, 0x0cfebacd,
				0x01188ecd, 0x06a96f84, 0x0ec138b4, 0x0a33392e,
				0x08f696a3, 0x09dc05b3, 0x0eeb3c87, 0x073bf2fd,
				0x07c41931, 0x0b66730f, 0x03950403, 0x05a33abb,
			},
			[16]uint32{
				0x02779e92, 0x069da0da, 0x059c7500, 0x02081cb7,
				0x024c1ab5, 0x06a4f24d, 0x016e8bf5, 0x01b73259,
				0x03278397, 0x063fbf98, 0x0f6a624d, 0x00c3f878,
				0x0aa90dbb, 0x01ec2b67, 0x0ed7b66d, 0x089f76c5,
			},
			[16]uint32{
				0x0b57922d, 0x08ea8a8a, 0x07751bc6, 0x02539061,
				0x03900a12, 0x057765e2, 0x093495f3, 0x0354b707,
				0x027ed773, 0x0fd1b87c, 0x05690b42, 0x048713a6,
				0x048cbd5b, 0x05db9f26, 0x0402b192, 0x00ad3145,
			},
			[16]uint32{
				0x00d32c88, 0x0473b737, 0x0ce7ece7, 0x0a623ab8,
				0x00a5f9b2, 0x0708f701, 0x0409da36, 0x077e8909,
				0x09178d8b, 0x014130c6, 0x05c5eaca, 0x0fb9e240,
				0x003f47ac, 0x0aa47b98, 0x08c891f6, 0x04329973,
			},
		),
		ed448.NewPoint(
			[16]uint32{
				0x028f0b14, 0x0c4e8cac, 0x01e6505f, 0x0e9387dd,
				0x0540ec47, 0x029dd2e5, 0x04d27275, 0x01b7753f,
				0x0ba58472, 0x01cd77c1, 0x0e9ee7a7, 0x0d00caeb,
				0x0a02a6eb, 0x00797cab, 0x01f1b3bb, 0x0672bccb,
			},
			[16]uint32{
				0x0896bbb2, 0x0e018c7f, 0x055337fa, 0x00cb0da2,
				0x0f59a101, 0x06453a63, 0x0d346171, 0x0efe2f17,
				0x0ea24201, 0x096a0446, 0x0f7eb766, 0x00270eb7,
				0x02c0fb78, 0x0b4cd897, 0x05cdb2e7, 0x0b670572,
			},
			[16]uint32{
				0x0d22a0cc, 0x02f89f49, 0x019e9f46, 0x083aaf9b,
				0x05a3a7be, 0x075421aa, 0x0a19dbcb, 0x05c515d2,
				0x0004ed88, 0x07c92fa3, 0x0feb6405, 0x0068b8de,
				0x06e58139, 0x00f59f55, 0x00d19953, 0x0f428949,
			},
			[16]uint32{
				0x0cd34494, 0x033ff9d9, 0x0ca379c2, 0x0b0f67df,
				0x098b2067, 0x06185880, 0x0680f73f, 0x017b76d9,
				0x07c1b26c, 0x0b72699b, 0x01019e9d, 0x097c7328,
				0x0dd32555, 0x010b1528, 0x0afe9352, 0x0555131a,
			},
		),
		ed448.NewPoint(
			[16]uint32{
				0x037325a7, 0x0929e858, 0x07a42aa4, 0x0a17e560,
				0x0a3dbf48, 0x096b7b9e, 0x0715946b, 0x0806de33,
				0x0492a576, 0x076a1f93, 0x01084d19, 0x0fde52e9,
				0x007c8644, 0x08a3f72a, 0x061bd87c, 0x00799405,
			},
			[16]uint32{
				0x007c3f7d, 0x08c1b95a, 0x036c668d, 0x02b90664,
				0x08ddc26e, 0x0b42bdb1, 0x021d0d82, 0x0fa3d577,
				0x0e11cd08, 0x000ca532, 0x071460ce, 0x0958d963,
				0x0c076a55, 0x0fb2d6ee, 0x055fc37d, 0x08b08940,
			},
			[16]uint32{
				0x0377e3e3, 0x079b22fd, 0x0217358e, 0x003242e6,
				0x0a65aa27, 0x0bceee30, 0x0238cf7b, 0x08bffd60,
				0x07d164d1, 0x00adbf86, 0x0b002060, 0x020c8896,
				0x0a0add70, 0x0cce917d, 0x09a2220c, 0x0f066560,
			},
			[16]uint32{
				0x0a74f626, 0x06e9c444, 0x028971ea, 0x0c3b0fb6,
				0x0f112cee, 0x0642f3f5, 0x050d9e74, 0x0c9b744d,
				0x0600a11d, 0x00bc86da, 0x0da32a37, 0x0245617b,
				0x042054c0, 0x0ff24cfc, 0x087373fb, 0x0ea1546d,
			},
		),
	}

	pub2 := &cramerShoupPublicKey{
		ed448.NewPoint(
			[16]uint32{
				0x02ad186a, 0x0ea87915, 0x0aa6b1e4, 0x064a5f2f,
				0x09c0d08d, 0x047ae943, 0x0eb1c2a8, 0x0e6769c2,
				0x0a6c88f1, 0x02f94c32, 0x01044cde, 0x028851a9,
				0x071c2398, 0x0b932bf1, 0x0d80e6eb, 0x05c3d697,
			},
			[16]uint32{
				0x0692af6f, 0x0d3a30e2, 0x0ea301de, 0x0dbddc31,
				0x0e01ede0, 0x0dc70521, 0x00a1a935, 0x00bd7618,
				0x0c353a9e, 0x0a278752, 0x036e75bf, 0x0afca072,
				0x079ae9db, 0x09d59130, 0x00a1ef5f, 0x0b96449f,
			},
			[16]uint32{
				0x00d45d75, 0x02a283ff, 0x057d717a, 0x00bb8508,
				0x0e611d83, 0x01450402, 0x0435abea, 0x0c7eff02,
				0x09bb69a1, 0x04f21f9e, 0x090a3122, 0x0260f1fa,
				0x082bcabe, 0x09ccf0a7, 0x01292cab, 0x0ebdcfb9,
			},
			[16]uint32{
				0x0f22d5b8, 0x0a62bf46, 0x078d5f81, 0x0445e577,
				0x06bb0c79, 0x03fe15c1, 0x010ac4ae, 0x021ac0bc,
				0x0a712e19, 0x074177eb, 0x02c3a252, 0x0894d930,
				0x0aab5528, 0x0d7b0b87, 0x072f0568, 0x0282bd46,
			},
		),
		ed448.NewPoint(
			[16]uint32{
				0x0cac7778, 0x037922e2, 0x0a5dc32d, 0x013172d6,
				0x0513804f, 0x0dac2fde, 0x02ceef98, 0x0a5ef879,
				0x0fc8dae2, 0x0141fb30, 0x0a7105f5, 0x03a52678,
				0x0c79e85f, 0x076fb526, 0x05205c2f, 0x0de36b3c,
			},
			[16]uint32{
				0x09bbb6aa, 0x0d40aa0a, 0x0ac92f3f, 0x020a75d4,
				0x0f624902, 0x098bc8cc, 0x06509c84, 0x0922a714,
				0x0613648c, 0x0dde0d93, 0x018246bf, 0x0f09b970,
				0x03d887d3, 0x04adc59f, 0x000a8b21, 0x0c81ab64,
			},
			[16]uint32{
				0x08ff730e, 0x0ed27055, 0x00ad18d6, 0x05130886,
				0x0409d513, 0x03132fdf, 0x0f1a9447, 0x016a9288,
				0x0cec8aa4, 0x07565d36, 0x07c3c27f, 0x07f2ce07,
				0x0bcacb80, 0x0efcc9dd, 0x04f2f578, 0x0a668ffd,
			},
			[16]uint32{
				0x043d5935, 0x04528c04, 0x0c57ff45, 0x00b00a5c,
				0x03d9189e, 0x02c13ac5, 0x0741a09c, 0x06a09b29,
				0x0d03d62e, 0x056c08ff, 0x057904c9, 0x051de63b,
				0x031916bd, 0x0ac016b6, 0x04f8f3dd, 0x0e7e2d1c,
			},
		),
		ed448.NewPoint(
			[16]uint32{
				0x08a85506, 0x062f9649, 0x0d75d751, 0x0142c488,
				0x039e6d2d, 0x09ab23b8, 0x0cd459bb, 0x0c91b8bc,
				0x099c2660, 0x0908f1fd, 0x0d7cc2fe, 0x0aa948ce,
				0x07c094b7, 0x02f08158, 0x0dff984c, 0x07845e02,
			},
			[16]uint32{
				0x0c33a54a, 0x058b50cf, 0x00c5cb5b, 0x03710be0,
				0x0d5382f7, 0x0293bdf1, 0x03d6acb6, 0x00cdfd25,
				0x09afc6f5, 0x00e6fa99, 0x0b62dbad, 0x05846556,
				0x070a5924, 0x09153c30, 0x005bf699, 0x067d2c52,
			},
			[16]uint32{
				0x0e1f1d60, 0x06a6f8d2, 0x011656e9, 0x059b3425,
				0x0793879f, 0x05064caf, 0x0e72f1b5, 0x05c184ec,
				0x0e948dbf, 0x072c1dec, 0x0b9e9a5f, 0x098cade9,
				0x0b662649, 0x04ece452, 0x03b2e4fd, 0x0fa3bd6a,
			},
			[16]uint32{
				0x022b7c53, 0x05a9ce25, 0x038d57f7, 0x0f7087c5,
				0x0456b5e5, 0x017cee64, 0x0b021582, 0x0abf7283,
				0x0df88f4e, 0x037306b1, 0x07b03628, 0x00aea180,
				0x0a962b1a, 0x04dbcfc8, 0x0605d4ea, 0x02825769,
			},
		),
	}

	message := []byte{
		0xd3, 0xff, 0x49, 0xad, 0x8b, 0x3e, 0x7a, 0x56,
		0x46, 0x9a, 0x7c, 0xe1, 0xdc, 0xb5, 0xb0, 0x2d,
		0x7a, 0x80, 0x82, 0x02, 0x4a, 0x99, 0x17, 0x30,
		0x0a, 0x36, 0x18, 0xcb, 0x17, 0xdb, 0xf4, 0x07,
		0xcc, 0x11, 0x7b, 0x36, 0xbf, 0x40, 0x41, 0x3c,
		0x9a, 0xa9, 0x5e, 0x07, 0xe8, 0xb9, 0x96, 0x72,
		0x91, 0x4a, 0x31, 0x92, 0x1c, 0xa0, 0xd6, 0x78,
	}

	expDRMessage := &drMessage{
		drCipher{
			// u11
			ed448.NewPoint(
				[16]uint32{
					0x0a17e677, 0x0c72738b, 0x044b91ea, 0x0aeee85a,
					0x097a1724, 0x0944b5a6, 0x0511f65a, 0x0e43351c,
					0x01012f88, 0x0afb7a50, 0x0b351695, 0x00db5165,
					0x0a4208f5, 0x046cb655, 0x04613291, 0x0a634ecc,
				},
				[16]uint32{
					0x0c79617c, 0x0c7fe24e, 0x01e77c1f, 0x0080bbf4,
					0x064227f8, 0x04c633e7, 0x074f4065, 0x0e4bffeb,
					0x08cdef09, 0x01f11868, 0x0819242d, 0x02dc8499,
					0x09544f24, 0x0b49eb65, 0x0a71cab9, 0x00562d85,
				},
				[16]uint32{
					0x005c962f, 0x08a14cc2, 0x09d9e8f5, 0x0d173a26,
					0x0b2be661, 0x036d1bc7, 0x0320d04d, 0x01cf0f71,
					0x07a2db99, 0x0ca1501c, 0x04db7eda, 0x0ffbf0c9,
					0x01850aa5, 0x0c2a6add, 0x0e15a4ec, 0x0a8e75e4,
				},
				[16]uint32{
					0x08480d1d, 0x0ef3dbcc, 0x037eb9f7, 0x03e74b03,
					0x09b0f66c, 0x043e7fb7, 0x00697125, 0x0378e43c,
					0x08740e77, 0x036135fb, 0x03930dca, 0x0964da2e,
					0x09b771fa, 0x00fef489, 0x0de92c8d, 0x0c35584d,
				},
			),
			// u21
			ed448.NewPoint(
				[16]uint32{
					0x0be81ed1, 0x07cd9250, 0x0f14c9a0, 0x00caddba,
					0x02c9f028, 0x0a8c516c, 0x0f5daf07, 0x0dda6081,
					0x008615e6, 0x00b47acf, 0x02cbfa2e, 0x0650e516,
					0x09ff4e42, 0x0e34bc86, 0x012aa6b4, 0x0da18443,
				},
				[16]uint32{
					0x01c66bf9, 0x05347f55, 0x070733bf, 0x083c3390,
					0x033f645f, 0x07f454a8, 0x020317ca, 0x09628235,
					0x0ea10625, 0x06772961, 0x0c132cbf, 0x0c93a6ae,
					0x0133b800, 0x09ad6c99, 0x0c3708c8, 0x0220d97a,
				},
				[16]uint32{
					0x0dd25779, 0x012fb2bd, 0x021ca196, 0x03f6288b,
					0x05915652, 0x01dbfecf, 0x01514d5c, 0x09ca77ff,
					0x0e8c3200, 0x012ed313, 0x082b0716, 0x0ab2d598,
					0x0882428f, 0x0e00f355, 0x0287d490, 0x0ef93f0a,
				},
				[16]uint32{
					0x076e2640, 0x094802ae, 0x0d7fd074, 0x091d9ef8,
					0x03b650f2, 0x0bddd3e5, 0x00d5f4e3, 0x02e0aa79,
					0x004bad45, 0x00a3b440, 0x09de886d, 0x067d3fc9,
					0x02b8223c, 0x093df563, 0x04ef4ca5, 0x0d82aefc,
				},
			),
			// e1
			ed448.NewPoint(
				[16]uint32{
					0x03264ddf, 0x0e7a3b2f, 0x006cc259, 0x07b5cedd,
					0x0acaaf96, 0x0a3ab6da, 0x0d59f8dc, 0x0735a85d,
					0x05d944fa, 0x03ef0547, 0x0f27d1e7, 0x0c66e93d,
					0x09024c34, 0x074b34df, 0x0e7df245, 0x0bb5a153,
				},
				[16]uint32{
					0x07c64cfc, 0x07c3e736, 0x0f29ec1d, 0x083cdb16,
					0x0deb70a6, 0x00a2755d, 0x06ec21df, 0x0edbba2d,
					0x0157a250, 0x05c845ea, 0x0e3ec6d2, 0x0dbd9da2,
					0x09f213cf, 0x09cd86ac, 0x00331554, 0x008385ed,
				},
				[16]uint32{
					0x04c38f6c, 0x021dbae2, 0x02d889fa, 0x0d6c8e80,
					0x00ea066a, 0x036351e2, 0x0fd34bbe, 0x0c0d48f6,
					0x0931c361, 0x085183e4, 0x080e2db9, 0x0ad70095,
					0x0fb3fc85, 0x03f5fb71, 0x0fdbfb73, 0x025be742,
				},
				[16]uint32{
					0x04a6ccd3, 0x06c62fe3, 0x01a821e0, 0x0b0a34f2,
					0x0d00aec7, 0x0c796592, 0x0d3f3c34, 0x0a822fc0,
					0x033476fa, 0x02ea07d2, 0x002e90e9, 0x0fe3bdc6,
					0x0cd599ce, 0x0255e7c9, 0x0dbf6de9, 0x07bd5ac8,
				},
			),
			// v1
			ed448.NewPoint(
				[16]uint32{
					0x06334d89, 0x05cd9c7c, 0x0b98ad75, 0x02d738ec,
					0x069eb49f, 0x0a50c661, 0x011a01cc, 0x0041588f,
					0x0f33f2f2, 0x0a675549, 0x0986a792, 0x01ad7c5f,
					0x01054ddb, 0x01f8b6fa, 0x0d72d85d, 0x045a6e77,
				},
				[16]uint32{
					0x0733522f, 0x0af4f4cb, 0x053afb69, 0x086135af,
					0x05521089, 0x0799da1d, 0x0bf34ca4, 0x0b3d04cf,
					0x0c005821, 0x0b24a051, 0x038a4840, 0x07745894,
					0x06351964, 0x0dfdbbcb, 0x0e3d04f2, 0x07f20357,
				},
				[16]uint32{
					0x050b7f7f, 0x0f0af399, 0x0dcb58ae, 0x0703f158,
					0x0e8e6e9f, 0x0315319c, 0x072f1ddf, 0x0f07f283,
					0x09ce1050, 0x0034bb83, 0x087dc8c7, 0x098872b1,
					0x0653708a, 0x05a5ed96, 0x076e9461, 0x0f7c8b58,
				},
				[16]uint32{
					0x008a5ea1, 0x0c15bf1a, 0x03678900, 0x074f6091,
					0x080e0116, 0x0bc51de5, 0x07a16525, 0x09ad79c7,
					0x0292418f, 0x0c42f824, 0x0153d038, 0x0986f73e,
					0x0e0bdbc6, 0x0ac54f18, 0x0723b299, 0x05784187,
				},
			),
			// u12
			ed448.NewPoint(
				[16]uint32{
					0x0c8b8396, 0x03aa75e7, 0x08c3dbbd, 0x018f2180,
					0x006904c9, 0x0c57fd7a, 0x075d14a6, 0x0504e045,
					0x04b6cf4d, 0x0fde7c99, 0x0ed1ed53, 0x096ab4cc,
					0x067993f2, 0x08d5cd2c, 0x0ce72cac, 0x0fba9428,
				},
				[16]uint32{
					0x0824d64e, 0x0e9783c1, 0x02e17d29, 0x0eec032e,
					0x0ff4b999, 0x0f4c526a, 0x00e44ded, 0x0d1915f4,
					0x0174c5c6, 0x07ad3d23, 0x04260041, 0x0944f671,
					0x005b695f, 0x06e26c1d, 0x0eea52d8, 0x030dd784,
				},
				[16]uint32{
					0x040e4131, 0x0f317b4d, 0x00a5f0c4, 0x00b0ffbd,
					0x0f02c1e3, 0x0e9512b1, 0x0ec742e3, 0x099f7b96,
					0x05542fcb, 0x05acd5fe, 0x02246935, 0x03b2dfdd,
					0x099a76ed, 0x0789fad9, 0x0addead7, 0x0985940c,
				},
				[16]uint32{
					0x0eb7d91c, 0x0abc62b8, 0x01607071, 0x0fdb6b8b,
					0x057d218a, 0x0f50e376, 0x04424d3e, 0x0080ab1a,
					0x012d00b9, 0x00187d84, 0x00cbfda6, 0x05684419,
					0x03d6444f, 0x04408f88, 0x02fd2e98, 0x0c00eaae,
				},
			),
			// u22
			ed448.NewPoint(
				[16]uint32{
					0x057e244f, 0x09842135, 0x07621c02, 0x053c7677,
					0x0b59f119, 0x02bd0778, 0x00946a29, 0x05fb8eba,
					0x02b9bcd1, 0x0cfffa34, 0x00fa277d, 0x06a77894,
					0x05898996, 0x050a7056, 0x0f4e5ba9, 0x02ca34fe,
				},
				[16]uint32{
					0x0894bb48, 0x06364bf6, 0x032bd738, 0x041b580d,
					0x08d7cc58, 0x0b4d8370, 0x05b32011, 0x03ecb176,
					0x0a7c79bf, 0x0f6f0b7c, 0x0a67356c, 0x02e3cf99,
					0x04f66417, 0x023de7e3, 0x06e2e74f, 0x0143841c,
				},
				[16]uint32{
					0x027e6abf, 0x0a146a3f, 0x02fa5fcb, 0x0f52285f,
					0x0e898ab3, 0x043d8f72, 0x077f99ab, 0x066ca58c,
					0x089391d7, 0x0f8e8a79, 0x01625814, 0x00735ff5,
					0x0e2c1e27, 0x03a5882c, 0x0efd15d4, 0x0e93c854,
				},
				[16]uint32{
					0x0752c266, 0x07baee88, 0x09b961dc, 0x073e0898,
					0x06a3f190, 0x0d16def6, 0x05c702d2, 0x01bb3ff9,
					0x0928c817, 0x0139fd2c, 0x0658862a, 0x02004992,
					0x0595d978, 0x030d4ecb, 0x0f5d93f3, 0x051490e8,
				},
			),
			// e2
			ed448.NewPoint(
				[16]uint32{
					0x0b99f3d2, 0x0802f9a5, 0x0b70dc0a, 0x03d17ea1,
					0x0ee1f47d, 0x0a8fe53f, 0x005f44fa, 0x072da748,
					0x05ba1a0c, 0x01f81b9a, 0x073ae1ca, 0x064917e0,
					0x0fec76f6, 0x0392a749, 0x0a9b2a18, 0x0d1699f7,
				},
				[16]uint32{
					0x02fc80c6, 0x0695fba7, 0x0fed28f0, 0x0eea1d26,
					0x09ee7484, 0x077ab819, 0x0e9e6333, 0x091a0e85,
					0x0329467b, 0x07f9bee4, 0x0039f161, 0x0020e9ba,
					0x06c57ace, 0x0dd78b68, 0x0c7d43eb, 0x0b1481d8,
				},
				[16]uint32{
					0x03c51aa3, 0x051c87fe, 0x06627098, 0x0d123c1d,
					0x09af6f19, 0x0fc7ef8c, 0x05775b03, 0x04141bbf,
					0x07f1593a, 0x0ad4a37f, 0x04fc0874, 0x058d9606,
					0x051763e7, 0x0825ca51, 0x0aae725c, 0x047b6fc9,
				},
				[16]uint32{
					0x04a4021e, 0x0f8d0cbe, 0x09190917, 0x079b2622,
					0x04d2e1a1, 0x0f6075fa, 0x0b6d9af7, 0x052f1e2b,
					0x01efba41, 0x08464680, 0x087d23e8, 0x09a5aa9b,
					0x00f2dfc7, 0x0f8d1425, 0x031c37a9, 0x0774be4e,
				},
			),
			// v2
			ed448.NewPoint(
				[16]uint32{
					0x091e6e57, 0x018646f1, 0x0d3d06fa, 0x06e3daa5,
					0x01634b17, 0x06aa68c8, 0x052d5ab7, 0x0cd7dfef,
					0x07f46c8b, 0x078271e7, 0x04fb394e, 0x084332c2,
					0x01d33251, 0x062696a2, 0x0e6fe46f, 0x097feb9f,
				},
				[16]uint32{
					0x0b303acf, 0x09795b60, 0x0106fc3e, 0x078a861a,
					0x05082369, 0x00d5ee11, 0x0a879854, 0x0b90992a,
					0x09fc2e23, 0x04681ba1, 0x049f3533, 0x0cee10f7,
					0x02018cc2, 0x0c5c40b1, 0x05e2a6b2, 0x0156878b,
				},
				[16]uint32{
					0x0f19abbc, 0x0bc15054, 0x0263f707, 0x02d9c8e1,
					0x0117ea48, 0x09e5047e, 0x0cb0fbcc, 0x0963e55e,
					0x0ee4d9a0, 0x0dd96439, 0x04e0327b, 0x0f7e45b1,
					0x071aefd9, 0x0881af15, 0x06d94ad6, 0x04b738f4,
				},
				[16]uint32{
					0x02af7386, 0x0de00412, 0x090d108d, 0x008f8222,
					0x02ca95ca, 0x070ae951, 0x08cea1b7, 0x0a3d9e8d,
					0x0db28902, 0x054a43a4, 0x01d9863a, 0x07ef0b9a,
					0x07c665b0, 0x0e1b58fa, 0x0a2198ac, 0x00ca674d,
				},
			),
		},
		nIZKProof{
			// l
			ed448.NewDecafScalar([]byte{
				0xbf, 0xaf, 0xc6, 0xc1, 0xc3, 0x74, 0xda, 0x23,
				0xa0, 0xd6, 0x92, 0xfe, 0x13, 0x39, 0xc7, 0xfd,
				0x17, 0xc4, 0xd4, 0x9e, 0x8d, 0xa2, 0x87, 0x89,
				0x34, 0x3b, 0x04, 0xde, 0xca, 0x95, 0x59, 0xbe,
				0x5c, 0x35, 0x29, 0xac, 0x55, 0x3c, 0xa0, 0xb7,
				0x7b, 0xff, 0xd3, 0xdd, 0x54, 0x0a, 0x2d, 0x3b,
				0x26, 0xfe, 0xfe, 0x4e, 0x94, 0x84, 0x2f, 0x0b,
			},
			),
			// n1
			ed448.NewDecafScalar([]byte{
				0xd8, 0x33, 0x30, 0x5f, 0xad, 0x6c, 0x50, 0x60,
				0x25, 0xa4, 0x31, 0x5c, 0xec, 0x31, 0xcb, 0xb0,
				0xe1, 0x13, 0xc5, 0xc8, 0x1b, 0x72, 0xe9, 0x4c,
				0xa8, 0xd3, 0x06, 0xb9, 0xc3, 0xdf, 0x95, 0xd1,
				0x83, 0x6e, 0x17, 0x14, 0x59, 0x19, 0x12, 0xea,
				0x51, 0x05, 0x1b, 0xbc, 0x5a, 0xb6, 0xb5, 0x0f,
				0xdc, 0xf8, 0x91, 0xcd, 0x15, 0x0e, 0xb5, 0x36,
			},
			),
			// n2
			ed448.NewDecafScalar([]byte{
				0xdd, 0xf5, 0x2e, 0x4b, 0x8a, 0x7f, 0x32, 0x83,
				0xd0, 0x7d, 0xd0, 0x31, 0x22, 0x86, 0x37, 0x9a,
				0x75, 0x48, 0x3b, 0x83, 0xb3, 0x02, 0xfd, 0x96,
				0x48, 0xbe, 0xfe, 0x65, 0xc1, 0x53, 0xa1, 0x5e,
				0x57, 0x8d, 0xc0, 0x87, 0x1a, 0x68, 0xcd, 0xfa,
				0xc2, 0x0e, 0x5e, 0x2e, 0x5b, 0x84, 0xef, 0xcc,
				0x79, 0xc2, 0x83, 0x13, 0x2b, 0x90, 0xec, 0x0d,
			},
			),
		},
	}

	drMessage := &drMessage{}
	err := drMessage.drEnc(message, fixedRand(randDREData), pub1, pub2)

	c.Assert(drMessage.cipher, DeepEquals, expDRMessage.cipher)
	c.Assert(drMessage.proof, DeepEquals, expDRMessage.proof)
	c.Assert(err, IsNil)
}

func (s *OTR4Suite) Test_DREncryptAndDecrypt(c *C) {
	message := []byte{
		0xfd, 0xf1, 0x18, 0xbf, 0x8e, 0xc9, 0x64, 0xc7,
		0x94, 0x46, 0x49, 0xda, 0xcd, 0xac, 0x2c, 0xff,
		0x72, 0x5e, 0xb7, 0x61, 0x46, 0xf1, 0x93, 0xa6,
		0x70, 0x81, 0x64, 0x37, 0x7c, 0xec, 0x6c, 0xe5,
		0xc6, 0x8d, 0x8f, 0xa0, 0x43, 0x23, 0x45, 0x33,
		0x73, 0x79, 0xa6, 0x48, 0x57, 0xbb, 0x0f, 0x70,
		0x63, 0x8c, 0x62, 0x26, 0x9e, 0x17, 0x5d, 0x22,
	}

	priv1, pub1, err := deriveCramerShoupKeys(rand.Reader)
	priv2, pub2, err := deriveCramerShoupKeys(rand.Reader)

	drMessage := &drMessage{}
	err = drMessage.drEnc(message, rand.Reader, pub1, pub2)

	expMessage1, err := drMessage.drDec(pub1, pub2, priv1, 1)
	expMessage2, err := drMessage.drDec(pub1, pub2, priv2, 2)
	c.Assert(expMessage1, DeepEquals, message)
	c.Assert(expMessage2, DeepEquals, message)
	c.Assert(err, IsNil)
}
