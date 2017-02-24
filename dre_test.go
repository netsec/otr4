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
	out, err := auth(fixedRand(randData), testPubA, testPubB, testPubC, testSec, message)

	c.Assert(out, DeepEquals, testSigma)
	c.Assert(err, IsNil)

	r := make([]byte, 56*5-1)
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
	pub1 := &cramerShoupPublicKey{
		ed448.NewPoint(
			[16]uint32{
				0x0ad46e4f, 0x09d5caab, 0x000af852, 0x0d5612b6,
				0x07a30f2a, 0x0865ec16, 0x0ed39a46, 0x05f2d18d,
				0x03674be7, 0x05b9016d, 0x00445209, 0x00c25901,
				0x0c7ef035, 0x0c95d3bf, 0x069ed7ef, 0x00e751bf,
			},
			[16]uint32{
				0x0e025c56, 0x09c31b30, 0x055c1f9e, 0x0383581a,
				0x0b6da69d, 0x05b612c4, 0x066dd8f1, 0x044868f5,
				0x0be2a63b, 0x00198b93, 0x08e2bdec, 0x03638823,
				0x02f26258, 0x0952385c, 0x0b8ef9b7, 0x09f632f0,
			},
			[16]uint32{
				0x0a1d2cc4, 0x011ddaf4, 0x012947d3, 0x0b209cd5,
				0x0dda011b, 0x0e5dd15d, 0x040de557, 0x0b72af34,
				0x072b5ce6, 0x031d1f38, 0x00216d9c, 0x00ec17d9,
				0x0440a02d, 0x04c86a7d, 0x099d46a6, 0x03c73164,
			},
			[16]uint32{
				0x0bfff232, 0x00346edb, 0x0fbeae3f, 0x04f34a79,
				0x0a769735, 0x02540ef8, 0x02ef50b6, 0x0161fc36,
				0x02153316, 0x02efe7a4, 0x0d3d8dd8, 0x036f28a3,
				0x02f0b34d, 0x078ccb25, 0x03ecad39, 0x065cefb1,
			},
		),
		ed448.NewPoint(
			[16]uint32{
				0x03387e3c, 0x04b30383, 0x00f7deef, 0x0b300e5f,
				0x06341b1a, 0x026f1691, 0x0d7d0b56, 0x0ea44438,
				0x0a94ade2, 0x05971b55, 0x0bd874b9, 0x0a8d2a56,
				0x02ccc4ed, 0x010525cc, 0x0ed2506c, 0x0e9be7cf,
			},
			[16]uint32{
				0x03e65dde, 0x03e5fc44, 0x065b94b0, 0x08e289fd,
				0x04c43927, 0x06230dfb, 0x03d08abe, 0x0273297e,
				0x0c3d07b8, 0x048ad9f6, 0x0f21566b, 0x06cdee16,
				0x05b81dac, 0x09573b8d, 0x07785b4f, 0x038142bc,
			},
			[16]uint32{
				0x09d9fc03, 0x09f2fe2e, 0x005abcd3, 0x0482cc8e,
				0x0bf189e2, 0x0884cd84, 0x08c0bc4f, 0x0d847aee,
				0x0bbcbc93, 0x08ce0af3, 0x0ca7d5b0, 0x025e4539,
				0x0653691d, 0x0f0f369d, 0x09d0aa32, 0x0f780bd9,
			},
			[16]uint32{
				0x084fed4c, 0x0e970c2e, 0x050716be, 0x0355f127,
				0x0b96da42, 0x023a0fa1, 0x0eeb8cdb, 0x061f87ca,
				0x022c2a49, 0x0aeb6b5c, 0x0982b2d4, 0x0b5c074c,
				0x0f44376b, 0x0f4b11a1, 0x052a053b, 0x04327cdf,
			},
		),
		ed448.NewPoint(
			[16]uint32{
				0x0abaefe4, 0x01710e33, 0x073d78fe, 0x08461452,
				0x06812257, 0x0232f36c, 0x0dbb29c6, 0x0d2db4f5,
				0x00e65088, 0x0f575ac4, 0x0b825170, 0x07dc4e48,
				0x0eef1437, 0x06755116, 0x06a8a6ff, 0x041d1ccc,
			},
			[16]uint32{
				0x00593959, 0x0c478000, 0x0f653465, 0x0923d3cf,
				0x028e910b, 0x05d49af0, 0x0651be3f, 0x0b074939,
				0x08f60286, 0x029ec003, 0x0c34f890, 0x0ff0fdb6,
				0x0cc41b62, 0x059de840, 0x032c8d6d, 0x0a2bc724,
			},
			[16]uint32{
				0x055ae1ba, 0x0c447604, 0x094a27f5, 0x0ed70f37,
				0x06523fa2, 0x09d6665d, 0x0bb2ce14, 0x024eb41a,
				0x07558290, 0x05f68265, 0x05c13e4c, 0x09adb060,
				0x0d8e7aae, 0x03287a3c, 0x08b2fc0e, 0x0f27cabd,
			},
			[16]uint32{
				0x095187ea, 0x0d9186c9, 0x05fc4c1e, 0x06e99b34,
				0x06d8c6a1, 0x0e5ccdc5, 0x00164428, 0x005a3020,
				0x05865c76, 0x0953ea61, 0x0dadf5f4, 0x0c300af1,
				0x09422e92, 0x0b0a0dbd, 0x097001c8, 0x0161ccdf,
			},
		),
	}

	pub2 := &cramerShoupPublicKey{
		ed448.NewPoint(
			[16]uint32{
				0x0ecc1c06, 0x02cd26a7, 0x03b559f4, 0x0ade3ef9,
				0x0a5a0a51, 0x001d442b, 0x0a3b41d7, 0x09df4c68,
				0x0a9a130c, 0x0c66c70e, 0x028438da, 0x01a618ac,
				0x096add53, 0x01a37c84, 0x06da6c38, 0x052f228b,
			},
			[16]uint32{
				0x0bba2e09, 0x081ffffe, 0x0c2b616f, 0x05b25369,
				0x0f00e105, 0x0917460b, 0x0596dfd7, 0x023a1b48,
				0x0db660b3, 0x076c145f, 0x00f78e70, 0x006e723e,
				0x04a3a0d4, 0x05a71853, 0x0da580f2, 0x0e6d3bc4,
			},
			[16]uint32{
				0x0cd512b1, 0x09d3ab57, 0x0c002f1d, 0x0d3ae17b,
				0x0c1d9737, 0x0f85498b, 0x0ea469e4, 0x09858571,
				0x068e06c8, 0x018d5ad5, 0x06e48fcb, 0x00172f60,
				0x0b42b6cb, 0x0147531a, 0x0a596011, 0x052ecaad,
			},
			[16]uint32{
				0x0a24e9b9, 0x0413798a, 0x03d1bf2b, 0x03545838,
				0x00cdd081, 0x0bbe864d, 0x0481e3d8, 0x03b7eb31,
				0x0cc7a3ae, 0x08abb739, 0x00a6ba9b, 0x00d08cb4,
				0x0f859b4d, 0x0e469483, 0x0e83aaa3, 0x007c5d34,
			},
		),
		ed448.NewPoint(
			[16]uint32{
				0x0c5f2c8f, 0x01fde2c9, 0x0713e6c4, 0x0ac2f953,
				0x00d36d8f, 0x0ed1574e, 0x02b1ff2c, 0x079c2dc3,
				0x0cd8b853, 0x03200803, 0x084d69f6, 0x0012dbc6,
				0x0a5c7a5f, 0x079e2057, 0x04a9f3cc, 0x0ac85481,
			},
			[16]uint32{
				0x089993a5, 0x0b6fe8bb, 0x09208db0, 0x02921a74,
				0x0e901eb5, 0x06a4df1d, 0x0c2e73d1, 0x06a4eed3,
				0x0b391352, 0x0b139591, 0x037484b9, 0x07f1bb24,
				0x09121b3f, 0x058c8c4f, 0x08136aa1, 0x06d99920,
			},
			[16]uint32{
				0x0ab7ecb5, 0x0e59eb9d, 0x0525bbd3, 0x029b6487,
				0x08429a1f, 0x0f0910d2, 0x05f3d665, 0x05be0edd,
				0x06dadf49, 0x06007c49, 0x0113a0d7, 0x0c696a3a,
				0x00760419, 0x0f5ec587, 0x06716bb1, 0x0e393abd,
			},
			[16]uint32{
				0x01369ba0, 0x044d5ab5, 0x0f6e6715, 0x01c5d42c,
				0x04b8a919, 0x0fc1f2a4, 0x053e13c0, 0x0bf4e206,
				0x03504df4, 0x01f0ce45, 0x0ceffed9, 0x056c265f,
				0x0e0a21e1, 0x066203d5, 0x07628b72, 0x0475582d,
			},
		),
		ed448.NewPoint(
			[16]uint32{
				0x0d82ae04, 0x07a78abd, 0x078cde78, 0x0e4edf72,
				0x0ed185bc, 0x0e32c186, 0x052bc5a6, 0x0c454c79,
				0x04cb6bf5, 0x042e4609, 0x0d094a7e, 0x0441ba1b,
				0x02dfc41d, 0x03775c13, 0x080455e4, 0x07e0c89c,
			},
			[16]uint32{
				0x0e278c68, 0x0892b8a9, 0x0e35bbad, 0x0926ff4c,
				0x03f9fbc5, 0x0686586a, 0x04675239, 0x08bfacf2,
				0x0474df13, 0x00072767, 0x056b78e8, 0x02d6a981,
				0x0e0025dc, 0x02f9c999, 0x064677ed, 0x0ffa8bba,
			},
			[16]uint32{
				0x0b4c6798, 0x05e7c037, 0x006be875, 0x0adceba2,
				0x0fc249db, 0x050119b9, 0x0b75f080, 0x0aceadd4,
				0x0d682fe6, 0x0665e297, 0x05466a1c, 0x09299454,
				0x0391bbd3, 0x0295bb92, 0x04f0d0fd, 0x0d6746c2,
			},
			[16]uint32{
				0x0150a7c1, 0x072a0266, 0x03094951, 0x01817d81,
				0x0a646a32, 0x0a7eb356, 0x01d8451a, 0x012c31e7,
				0x08e5461b, 0x03bb4265, 0x05efbe94, 0x0b798b68,
				0x0562166b, 0x0f6f24c2, 0x0d2c32a9, 0x071945d8,
			},
		),
	}

	message := []byte{
		0x80, 0x47, 0xa4, 0x8b, 0xa7, 0xa5, 0x94, 0x58,
		0xee, 0x80, 0x0f, 0xd1, 0xc8, 0x7f, 0x7a, 0xae,
		0x4c, 0xb1, 0xa6, 0xdd, 0x07, 0xa6, 0xc5, 0xda,
		0x59, 0xb7, 0x54, 0x84, 0x76, 0xb1, 0x48, 0xef,
		0x21, 0x59, 0xba, 0xfe, 0x8d, 0x02, 0x16, 0x53,
		0xf0, 0x4e, 0xb6, 0x23, 0xef, 0x03, 0x4c, 0x7c,
		0x59, 0x39, 0xd5, 0x43, 0xda, 0xed, 0x28, 0x62,
	}

	randData := []byte{
		0xc9, 0x21, 0xa6, 0x41, 0xc3, 0x43, 0xb3, 0x4f,
		0x3e, 0x86, 0x99, 0xbf, 0x11, 0x75, 0x2c, 0x40,
		0x05, 0xb9, 0x0e, 0xd1, 0x01, 0xd8, 0x3e, 0xeb,
		0xda, 0xfa, 0x7e, 0x28, 0x94, 0xe8, 0x62, 0x31,
		0xa5, 0x62, 0xfd, 0x27, 0x85, 0x00, 0xdf, 0x4a,
		0xc3, 0xc2, 0x27, 0x2e, 0x11, 0x49, 0xfc, 0x3c,
		0xc0, 0xdf, 0x80, 0x3d, 0x7a, 0x2f, 0x1f, 0x06,
		0xc9, 0x21, 0xa6, 0x41, 0xc3, 0x43, 0xb3, 0x4f,
		0x3e, 0x86, 0x99, 0xbf, 0x11, 0x75, 0x2c, 0x40,
		0x05, 0xb9, 0xff, 0xd1, 0x01, 0xd8, 0x3e, 0xeb,
		0xda, 0xfa, 0x7e, 0x28, 0x20, 0xe8, 0x62, 0x31,
		0xa5, 0x34, 0xfd, 0x27, 0x85, 0x00, 0xdd, 0x4a,
		0xcc, 0xc2, 0x27, 0xee, 0x11, 0x10, 0xfc, 0x3c,
		0xc0, 0xdf, 0x80, 0x3d, 0x7a, 0x2f, 0x1f, 0x06,
	}

	expV1 := ed448.NewPoint(
		[16]uint32{
			0x0f0c83b1, 0x0081c017, 0x0722baba, 0x0a9f20ca,
			0x0ab28cf0, 0x0caf8aec, 0x0992a0f7, 0x0187d7a2,
			0x0f0d7981, 0x0c38d08d, 0x0e610473, 0x0fc52752,
			0x05c78621, 0x01095896, 0x03ff82b1, 0x0c710192,
		},
		[16]uint32{
			0x0d536906, 0x09675e67, 0x0e5fcc06, 0x04d6042d,
			0x058ae2e2, 0x0ba51e83, 0x0b113bbb, 0x0491dbd9,
			0x017d988e, 0x0c2bf7ee, 0x07c8a6b1, 0x0e294703,
			0x00a24910, 0x0b12ebf5, 0x0565cd97, 0x06a20cb0,
		},
		[16]uint32{
			0x0f58107a, 0x0a0b63a7, 0x009021a4, 0x0672ed2e,
			0x0d4466b5, 0x086d411c, 0x0f023ad6, 0x0a25b4f3,
			0x0efc1c5e, 0x0c223917, 0x05ce7e7d, 0x0503246c,
			0x0033d84e, 0x00049394, 0x0f112383, 0x04f7ff30,
		},
		[16]uint32{
			0x02966f2b, 0x037bac22, 0x04f6ffe3, 0x0786e355,
			0x0dc720ac, 0x0d6b3fa0, 0x019c2716, 0x0a5295fe,
			0x018f81a4, 0x06d8c231, 0x082d8714, 0x0e3de239,
			0x0379a9a5, 0x035f03e6, 0x03b301fe, 0x0c3e67d5,
		},
	)

	expV2 := ed448.NewPoint(
		[16]uint32{
			0x032f40d3, 0x0e2d0aa6, 0x0cfa7ce2, 0x058b6e92,
			0x08080d13, 0x0b97a933, 0x0a63ad00, 0x020d9604,
			0x0a59a858, 0x0a34047a, 0x066639ea, 0x0b2c0d75,
			0x0ad0f40e, 0x03d9f6a4, 0x01eb5588, 0x0d62332b,
		},
		[16]uint32{
			0x0d8bf918, 0x056a6052, 0x0a09c330, 0x08e290d5,
			0x0a34fcc7, 0x09c9c52c, 0x09424048, 0x082ad610,
			0x05911b71, 0x016d700a, 0x0828793e, 0x05a227cf,
			0x01f509c0, 0x0435367e, 0x0acefc4f, 0x0a4fef98,
		},
		[16]uint32{
			0x05477c87, 0x06ad7449, 0x00a9bb27, 0x05a1db91,
			0x0607e107, 0x06a97568, 0x012286f7, 0x07e9b251,
			0x00e7f41d, 0x092a88ef, 0x039f9890, 0x06110778,
			0x0e921bfd, 0x0cadb6ff, 0x0942c869, 0x04a990d4,
		},
		[16]uint32{
			0x076a864b, 0x04e9af4b, 0x09e6edbf, 0x01135ee6,
			0x03552ff0, 0x0659d362, 0x061b2a64, 0x0d490b87,
			0x0c357d5d, 0x0a787e7e, 0x078c797b, 0x0ac189c4,
			0x0b65d079, 0x071cdcb1, 0x0ff192b9, 0x034eb9f3,
		},
	)

	v1, v2, err := drEnc(message, fixedRand(randData), pub1, pub2)
	c.Assert(v1, DeepEquals, expV1)
	c.Assert(v2, DeepEquals, expV2)
	c.Assert(err, IsNil)
}
