package core

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"

	emath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/vars"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var yuckyGlobalTestEnableMess = false

func runMESSTest(t *testing.T, easyL, hardL, caN int, easyT, hardT int64) (hardHead bool, err error) {
	// Generate the original common chain segment and the two competing forks
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	genesisB := MustCommitGenesis(db, genesis)

	chain, err := NewBlockChain(db, nil, genesis.Config, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(yuckyGlobalTestEnableMess)

	easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, easyL, func(i int, b *BlockGen) {
		b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
		b.OffsetTime(easyT)
	})
	commonAncestor := easy[caN-1]
	hard, _ := GenerateChain(genesis.Config, commonAncestor, engine, db, hardL, func(i int, b *BlockGen) {
		b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
		b.OffsetTime(hardT)
	})

	if _, err := chain.InsertChain(easy); err != nil {
		t.Fatal(err)
	}
	_, err = chain.InsertChain(hard)
	hardHead = chain.CurrentBlock().Hash() == hard[len(hard)-1].Hash()
	return
}

func TestBlockChain_AF_ECBP1100(t *testing.T) {
	t.Skip("These have been disused as of the sinusoidal -> cubic change.")
	yuckyGlobalTestEnableMess = true
	defer func() {
		yuckyGlobalTestEnableMess = false
	}()

	cases := []struct {
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset            int64
		hardGetsHead, accepted            bool
	}{
		// INDEX=0
		// Hard has insufficient total difficulty / length and is rejected.
		{
			5000, 7500, 2500,
			50, -9,
			false, false,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			1000, 7, 995,
			60, 0,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			1000, 7, 995,
			60, 7,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			1000, 1, 999,
			30, 1,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 3, 497,
			0, -8,
			true, true,
		},
		// INDEX=5
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 4, 496,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 5, 495,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 6, 494,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 7, 493,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 8, 492,
			0, -9,
			true, true,
		},
		// INDEX=10
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 9, 491,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 12, 488,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 20, 480,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 40, 460,
			0, -9,
			true, true,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 60, 440,
			0, -9,
			true, true,
		},
		// // INDEX=15
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 250, 250,
			0, -9,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 250, 250,
			7, -9,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 300, 200,
			13, -9,
			false, false,
		},
		// Hard has sufficient total difficulty / length and is accepted.
		{
			500, 200, 300,
			47, -9,
			true, true,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			47, -8,
			false, false,
		},
		// // INDEX=20
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			17, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			7, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 200, 300,
			0, -8,
			false, false,
		},
		// Hard has insufficient total difficulty / length and is rejected.
		{
			500, 100, 400,
			0, -7,
			false, false,
		},
		// Hard is accepted, but does not have greater total difficulty,
		// and is not set as the chain head.
		{
			1000, 1, 900,
			60, -9,
			false, true,
		},
		// INDEX=25
		// Hard is shorter, but sufficiently heavier chain, is accepted.
		{
			500, 100, 390,
			60, -9,
			true, true,
		},
	}

	for i, c := range cases {
		hardHead, err := runMESSTest(t, c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset)
		if (err != nil && c.accepted) || (err == nil && !c.accepted) || (hardHead != c.hardGetsHead) {
			t.Errorf("case=%d [easy=%d hard=%d ca=%d eo=%d ho=%d] want.accepted=%v want.hardHead=%v got.hardHead=%v err=%v",
				i,
				c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset,
				c.accepted, c.hardGetsHead, hardHead, err)
		}
	}
}

func TestBlockChain_AF_ECBP1100_2(t *testing.T) {
	yuckyGlobalTestEnableMess = true
	defer func() {
		yuckyGlobalTestEnableMess = false
	}()

	cases := []struct {
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset            int64
		hardGetsHead, accepted            bool
	}{
		// Random coin tosses involved for equivalent difficulty.
		// {
		// 	1000, 1, 999,
		// 	0, 0, // -1 offset => 10-1=9 same child difficulty
		// 	false, true,
		// },
		// {
		// 	1000, 3, 997,
		// 	0, 0, // -1 offset => 10-1=9 same child difficulty
		// 	false, true,
		// },
		// {
		// 	1000, 10, 990,
		// 	0, 0, // -1 offset => 10-1=9 same child difficulty
		// 	false, true,
		// },
		{
			1000, 1, 999,
			0, -2, // better difficulty
			true, true,
		},
		{
			1000, 25, 975,
			0, -2, // better difficulty
			true, true,
		},
		{
			1000, 30, 970,
			0, -2, // better difficulty
			false, true,
		},
		{
			1000, 50, 950,
			0, -5,
			true, true,
		},
		{
			1000, 50, 950,
			0, -1,
			false, true,
		},
		{
			1000, 999, 1,
			0, -9,
			true, true,
		},
		{
			1000, 999, 1,
			0, -8,
			false, true,
		},
		{
			1000, 500, 500,
			0, -8,
			true, true,
		},
		{
			1000, 500, 500,
			0, -7,
			false, true,
		},
		{
			1000, 300, 700,
			0, -7,
			false, true,
		},
		// Will pass, takes a long time.
		// {
		// 	5000, 4000, 1000,
		// 	0, -9,
		// 	true, true,
		// },
	}

	for i, c := range cases {
		hardHead, err := runMESSTest(t, c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset)
		if (err != nil && c.accepted) || (err == nil && !c.accepted) || (hardHead != c.hardGetsHead) {
			t.Errorf("case=%d [easy=%d hard=%d ca=%d eo=%d ho=%d] want.accepted=%v want.hardHead=%v got.hardHead=%v err=%v",
				i,
				c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset,
				c.accepted, c.hardGetsHead, hardHead, err)
		}
	}
}

func TestBlockChain_GenerateMESSPlot(t *testing.T) {
	// t.Skip("This test plots graph of chain acceptance for visualization.")

	easyLen := 500
	maxHardLen := 400

	generatePlot := func(title, fileName string) {
		p, err := plot.New()
		if err != nil {
			log.Panic(err)
		}
		p.Title.Text = title
		p.X.Label.Text = "Block Depth"
		p.Y.Label.Text = "Mode Block Time Offset (10 seconds + y)"

		accepteds := plotter.XYs{}
		rejecteds := plotter.XYs{}
		sides := plotter.XYs{}

		for i := 1; i <= maxHardLen; i++ {
			for j := -9; j <= 8; j++ {
				fmt.Println("running", i, j)
				hardHead, err := runMESSTest(t, easyLen, i, easyLen-i, 0, int64(j))
				point := plotter.XY{X: float64(i), Y: float64(j)}
				if err == nil && hardHead {
					accepteds = append(accepteds, point)
				} else if err == nil && !hardHead {
					sides = append(sides, point)
				} else if err != nil {
					rejecteds = append(rejecteds, point)
				}

				if err != nil {
					t.Log(err)
				}
			}
		}

		scatterAccept, _ := plotter.NewScatter(accepteds)
		scatterReject, _ := plotter.NewScatter(rejecteds)
		scatterSide, _ := plotter.NewScatter(sides)

		pixelWidth := vg.Length(1000)

		scatterAccept.Color = color.RGBA{R: 152, G: 236, B: 161, A: 255}
		scatterAccept.Shape = draw.BoxGlyph{}
		scatterAccept.Radius = vg.Length((float64(pixelWidth) / float64(maxHardLen)) * 2 / 3)
		scatterReject.Color = color.RGBA{R: 236, G: 106, B: 94, A: 255}
		scatterReject.Shape = draw.BoxGlyph{}
		scatterReject.Radius = vg.Length((float64(pixelWidth) / float64(maxHardLen)) * 2 / 3)
		scatterSide.Color = color.RGBA{R: 190, G: 197, B: 236, A: 255}
		scatterSide.Shape = draw.BoxGlyph{}
		scatterSide.Radius = vg.Length((float64(pixelWidth) / float64(maxHardLen)) * 2 / 3)

		p.Add(scatterAccept)
		p.Legend.Add("Accepted", scatterAccept)
		p.Add(scatterReject)
		p.Legend.Add("Rejected", scatterReject)
		p.Add(scatterSide)
		p.Legend.Add("Sidechained", scatterSide)

		p.Legend.YOffs = -30

		err = p.Save(pixelWidth, 300, fileName)
		if err != nil {
			log.Panic(err)
		}
	}
	yuckyGlobalTestEnableMess = true
	defer func() {
		yuckyGlobalTestEnableMess = false
	}()
	baseTitle := fmt.Sprintf("Accept/Reject Reorgs: Relative Time (Difficulty) over Proposed Segment Length (%d-block original chain)", easyLen)
	generatePlot(baseTitle, "reorgs-MESS.png")
	yuckyGlobalTestEnableMess = false
	// generatePlot("WITHOUT MESS: "+baseTitle, "reorgs-noMESS.png")
}

// Some weird constants to avoid constant memory allocs for them.
var (
	bigMinus99 = big.NewInt(-99)
	big1       = big.NewInt(1)
)

func TestBlockChain_GenerateMESSPlot_2(t *testing.T) {
	// t.Skip("This test plots graph of chain acceptance for visualization.")

	localBlockTime := uint64(10)
	baseTotalSpanSeconds := localBlockTime * 1000        // ie 1000 blocks at 10 seconds each
	reorgSpanSecondsMax := baseTotalSpanSeconds * 9 / 10 // how far back our proposed reorg is going to go

	baseBlockLen := baseTotalSpanSeconds / localBlockTime

	segmentTotalDifficulty := func(segment []*big.Int) *big.Int {
		out := big.NewInt(0)
		for _, b := range segment {
			out.Add(out, b)
		}
		return out
	}

	generatePlot := func(title, fileName string, useProposedSpan bool) {
		p, err := plot.New()
		if err != nil {
			log.Panic(err)
		}
		p.Title.Text = title
		p.Y.Label.Text = "Block Depth"
		p.X.Label.Text = "Mode Block Time Offset (10 seconds + y)"

		accepteds := plotter.XYs{}
		rejecteds := plotter.XYs{}
		sides := plotter.XYs{}


		// - duration: how long the chain should span in seconds
		// - maxDifficulty: maximimum difficulty allowed per block (eg as relative hashrate)
		generateDifficultySet := func(initDifficulty, maxDifficulty *big.Int, duration int64) []*big.Int {

			if maxDifficulty == nil {
				maxDifficulty = big.NewInt(math.MaxInt64)
			}

			parentDiff := new(big.Int).Set(initDifficulty)
			blockTime := uint64(1)

			outset := []*big.Int{}
			for duration > 0 {

				// parent difficulty met the max (allowed hashrate fulfilled)
				if parentDiff.Cmp(maxDifficulty) >= 0 {
					blockTime = 9
				}
				duration -= int64(blockTime)
				
				// https://github.com/ethereum/EIPs/issues/100
				// algorithm:
				// diff = (parent_diff +
				//         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
				//        ) + 2^(periodCount - 2)
				out := new(big.Int)
				out.Div(new(big.Int).SetUint64(blockTime), vars.EIP100FDifficultyIncrementDivisor)

				// if parent.UncleHash == types.EmptyUncleHash {
				// 	out.Sub(big1, out)
				// } else {
				// 	out.Sub(big2, out)
				// }
				out.Sub(big1, out)

				out.Set(emath.BigMax(out, bigMinus99))

				out.Mul(new(big.Int).Div(parentDiff, vars.DifficultyBoundDivisor), out)
				out.Add(out, parentDiff)

				// after adjustment and before bomb
				out.Set(emath.BigMax(out, vars.MinimumDifficulty))

				parentDiff.Set(out) // set for next iteration

				outset = append(outset, out)
			}
			return outset
		}

		easy := generateDifficultySet(
			params.DefaultMessNetGenesisBlock().Difficulty,
			nil,
			uint64(10),
		)

		for proposedSpanSeconds := uint64(1); proposedSpanSeconds <= reorgSpanSecondsMax; proposedSpanSeconds++ {
			for hardTime := uint64(1); hardTime <= 18; hardTime++ {
				// fmt.Println("running", proposedSpanSeconds, hardTime)

				// fitBlocks accounts for the number of blocks
				// produceable in a time span when you're making
				// them faster (10 blocks in 10 seconds with 1-second offsets,
				// vs 1 block in 10 seconds with 10-second offset)
				fitBlocks := easyOffset / hardTime // 10:1 => 10, 10:2 => 5
				if fitBlocks == 0 {
					fitBlocks = uint64(proposedSpanSeconds) // hard.len
				} else {
					fitBlocks = uint64(proposedSpanSeconds) / fitBlocks
				}

				hard := generateDifficultySet(
					// easy[len(easy)-proposedSpanSeconds], // offset = hard.len / fitBlocks
					easy[uint64(len(easy))-proposedSpanSeconds], // offset = hard.len / fitBlocks
					nil,
					hardTime,
				)

				// localTD := segmentTotalDifficulty(easy[len(easy)-proposedSpanSeconds:])
				// localTD := segmentTotalDifficulty(easy)

				overwrittenEasyBlocks := fitBlocks
				if overwrittenEasyBlocks > easyOffset *proposedSpanSeconds {
					overwrittenEasyBlocks = easyOffset * proposedSpanSeconds
				}

				localTD := segmentTotalDifficulty(easy[uint64(len(easy))-overwrittenEasyBlocks-1:])
				// localTD := segmentTotalDifficulty(easy[uint64(len(easy))-proposedSpanSeconds-1:])
				propTD := segmentTotalDifficulty(hard)

				localSpan := int64(overwrittenEasyBlocks * 10)
				xBig := big.NewInt(localSpan)
				if useProposedSpan {
					propSpan := int64(proposedSpanSeconds * hardTime)
					xBig = big.NewInt(propSpan)
				}

				eq := ecbp1100PolynomialV(xBig)

				want := eq.Mul(eq, localTD)

				got := new(big.Int).Mul(propTD, ecbp1100PolynomialVCurveFunctionDenominator)

				nogo := got.Cmp(want) < 0

				point := plotter.XY{Y: -float64(proposedSpanSeconds), X: float64(hardTime)}

				reorg := propTD.Cmp(localTD) > 0
				if localTD.Cmp(propTD) == 0 && proposedSpanSeconds == overwrittenEasyBlocks && rand.Float64() < 0.5 {
					reorg = true
				}

				if reorg {
					if nogo {
						reorg = false
					}
				}

				// ok := eyalSirer || !nogo
				if nogo {
					// Would be side chain
					sides = append(sides, point)
				} else {
					accepteds = append(accepteds, point)
				}
			}
		}

		scatterAccept, _ := plotter.NewScatter(accepteds)
		scatterReject, _ := plotter.NewScatter(rejecteds)
		scatterSide, _ := plotter.NewScatter(sides)

		pixelWidth := vg.Length(reorgSpanSecondsMax * 110 / 100)

		scatterAccept.Color = color.RGBA{R: 0, G: 200, B: 11, A: 255}
		scatterAccept.Shape = draw.BoxGlyph{}
		scatterAccept.Radius = vg.Length((float64(pixelWidth) / float64(reorgSpanSecondsMax)) * 2 / 3)
		scatterReject.Color = color.RGBA{R: 236, G: 106, B: 94, A: 255}
		scatterReject.Shape = draw.BoxGlyph{}
		scatterReject.Radius = vg.Length((float64(pixelWidth) / float64(reorgSpanSecondsMax)) * 2 / 3)
		scatterSide.Color = color.RGBA{R: 220, G: 227, B: 246, A: 255}
		scatterSide.Shape = draw.BoxGlyph{}
		scatterSide.Radius = vg.Length((float64(pixelWidth) / float64(reorgSpanSecondsMax)) * 2 / 3)

		p.Add(scatterAccept)
		p.Add(scatterReject)
		p.Add(scatterSide)

		// p.Legend.Add("Accepted", scatterAccept)
		// p.Legend.Add("Rejected", scatterReject)
		// p.Legend.Add("Sidechained", scatterSide)
		// p.Legend.YOffs = -30

		err = p.Save(600 ,pixelWidth, fileName)
		if err != nil {
			log.Panic(err)
		}
	}
	yuckyGlobalTestEnableMess = true
	defer func() {
		yuckyGlobalTestEnableMess = false
	}()
	baseTitle := fmt.Sprintf("Accept/Reject Reorgs: Relative Time (Difficulty) over Proposed Segment Length (%d-block original chain)", baseBlockLen)
	generatePlot(baseTitle, "reorgs-fast-localspan-MESS.png", false)
	// generatePlot(baseTitle, "reorgs-fast-proposedspan-MESS.png", true)
	yuckyGlobalTestEnableMess = false
	// generatePlot("WITHOUT MESS: "+baseTitle, "reorgs-noMESS.png")
}

func TestEcbp1100AGSinusoidalA(t *testing.T) {
	cases := []struct {
		in, out float64
	}{
		{0, 1},
		{25132, 31},
	}
	tolerance := 0.0000001
	for i, c := range cases {
		if got := ecbp1100AGSinusoidalA(c.in); got < c.out-tolerance || got > c.out+tolerance {
			t.Fatalf("%d: in: %0.6f want: %0.6f got: %0.6f", i, c.in, c.out, got)
		}
	}
}

/*
TestAFKnownBlock tests that AF functionality works for chain re-insertions.

Chain re-insertions use BlockChain.writeKnownBlockAsHead, where first-pass insertions
will hit writeBlockWithState.

AF needs to be implemented at both sites to prevent re-proposed chains from sidestepping
the AF criteria.
*/
func TestAFKnownBlock(t *testing.T) {
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	// genesis.Timestamp = 1
	genesisB := MustCommitGenesis(db, genesis)

	chain, err := NewBlockChain(db, nil, genesis.Config, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(true)

	easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, 1000, func(i int, gen *BlockGen) {
		gen.OffsetTime(0)
	})
	easyN, err := chain.InsertChain(easy)
	if err != nil {
		t.Fatal(err)
	}
	hard, _ := GenerateChain(genesis.Config, easy[easyN-300], engine, db, 300, func(i int, gen *BlockGen) {
		gen.OffsetTime(-7)
	})
	// writeBlockWithState
	if _, err := chain.InsertChain(hard); err != nil {
		t.Error("hard 1 not inserted (should be side)")
	}
	// writeKnownBlockAsHead
	if _, err := chain.InsertChain(hard); err != nil {
		t.Error("hard 2 inserted (will have 'ignored' known blocks, and never tried a reorg)")
	}
	hardHeadHash := hard[len(hard)-1].Hash()
	if chain.CurrentBlock().Hash() == hardHeadHash {
		t.Fatal("hard block got chain head, should be side")
	}
	if h := chain.GetHeaderByHash(hardHeadHash); h == nil {
		t.Fatal("missing hard block (should be imported as side, but still available)")
	}
}

func TestPlot_ecbp1100PolynomialV(t *testing.T) {
	t.Skip("This test plots a graph of the ECBP1100 polynomial curve.")
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "ECBP1100 Polynomial Curve Function"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	poly := plotter.NewFunction(func(f float64) float64 {
		n := big.NewInt(int64(f))
		y := ecbp1100PolynomialV(n)
		ff, _ := new(big.Float).SetInt(y).Float64()
		return ff
	})
	p.Add(poly)

	p.X.Min = 0
	p.X.Max = 30000
	p.Y.Min = 0
	p.Y.Max = 5000

	p.Y.Label.Text = "Antigravity imposition"
	p.X.Label.Text = "Seconds difference between local head and proposed common ancestor"

	if err := p.Save(1000, 1000, "ecbp1100-polynomial.png"); err != nil {
		t.Fatal(err)
	}
}

func TestEcbp1100PolynomialV(t *testing.T) {
	t.Log(
		ecbp1100PolynomialV(big.NewInt(99)),
		ecbp1100PolynomialV(big.NewInt(999)),
		ecbp1100PolynomialV(big.NewInt(99999)))
}

func TestGenerateChainTargetingHashrate(t *testing.T) {
	t.Skip("A development test to play with difficulty steps.")
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	// genesis.Timestamp = 1
	genesisB := MustCommitGenesis(db, genesis)

	chain, err := NewBlockChain(db, nil, genesis.Config, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(true)

	easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, 1000, func(i int, gen *BlockGen) {
		gen.OffsetTime(0)
	})
	if _, err := chain.InsertChain(easy); err != nil {
		t.Fatal(err)
	}
	firstDifficulty := chain.CurrentHeader().Difficulty
	targetDifficultyRatio := big.NewInt(2)
	targetDifficulty := new(big.Int).Div(firstDifficulty, targetDifficultyRatio)
	for chain.CurrentHeader().Difficulty.Cmp(targetDifficulty) > 0 {
		next, _ := GenerateChain(genesis.Config, chain.CurrentBlock(), engine, db, 1, func(i int, gen *BlockGen) {
			gen.OffsetTime(8) // 8: (=10+8=18>(13+4=17).. // minimum value over stable range
		})
		if _, err := chain.InsertChain(next); err != nil {
			t.Fatal(err)
		}
	}
	t.Log(chain.CurrentBlock().Number())
}

func TestBlockChain_AF_Difficulty_Develop(t *testing.T) {
	t.Skip("Development version of tests with plotter")
	// Generate the original common chain segment and the two competing forks
	engine := ethash.NewFaker()

	db := rawdb.NewMemoryDatabase()
	genesis := params.DefaultMessNetGenesisBlock()
	// genesis.Timestamp = 1
	genesisB := MustCommitGenesis(db, genesis)

	chain, err := NewBlockChain(db, nil, genesis.Config, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer chain.Stop()
	chain.EnableArtificialFinality(true)

	cases := []struct {
		easyLen, hardLen, commonAncestorN int
		easyOffset, hardOffset            int64
		hardGetsHead, accepted            bool
	}{
		// {
		// 	1000, 800, 200,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 800, 200,
		// 	60, 1,
		// 	true, true,
		// },
		// {
		// 	10000, 8000, 2000,
		// 	60, 1,
		// 	true, true,
		// },
		// {
		// 	20000, 18000, 2000,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	20000, 18000, 2000,
		// 	60, 1,
		// 	true, true,
		// },
		// {
		// 	10000, 8000, 2000,
		// 	10, 20,
		// 	true, true,
		// },

		// {
		// 	1000, 1, 999,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 10, 990,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 100, 900,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 200, 800,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 500, 500,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 999, 1,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	5000, 4000, 1000,
		// 	10, 1,
		// 	true, true,
		// },

		// {
		// 	10000, 9000, 1000,
		// 	10, 1,
		// 	true, true,
		// },
		//
		// {
		// 	7000, 6500, 500,
		// 	10, 1,
		// 	true, true,
		// },

		// {
		// 	100, 90, 10,
		// 	10, 1,
		// 	true, true,
		// },

		// {
		// 	1000, 1, 999,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 2, 998,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 3, 997,
		// 	10, 1,
		// 	true, true,
		// },
		// {
		// 	1000, 1, 999,
		// 	10, 8,
		// 	true, true,
		// },

		{
			1000, 50, 950,
			10, 9,
			false, false,
		},
		{
			1000, 100, 900,
			10, 8,
			false, false,
		},
		{
			1000, 100, 900,
			10, 7,
			false, false,
		},
		{
			1000, 50, 950,
			10, 5,
			true, true,
		},
		{
			1000, 50, 950,
			10, 3,
			true, true,
		},
		//5
		{
			1000, 100, 900,
			10, 3,
			false, false,
		},
		{
			1000, 200, 800,
			10, 3,
			false, false,
		},
		{
			1000, 200, 800,
			10, 1,
			false, false,
		},
	}

	// poissonTime := func(b *BlockGen, seconds int64) {
	// 	poisson := distuv.Poisson{Lambda: float64(seconds)}
	// 	r := poisson.Rand()
	// 	if r < 1 {
	// 		r = 1
	// 	}
	// 	if r > float64(seconds) * 1.5 {
	// 		r = float64(seconds)
	// 	}
	// 	chainreader := &fakeChainReader{config: b.config}
	// 	b.header.Time = b.parent.Time() + uint64(r)
	// 	b.header.Difficulty = b.engine.CalcDifficulty(chainreader, b.header.Time, b.parent.Header())
	// 	for err := b.engine.VerifyHeader(chainreader, b.header, false);
	// 		err != nil && err != consensus.ErrUnknownAncestor && b.header.Time > b.parent.Header().Time; {
	// 		t.Log(err)
	// 		r -= 1
	// 		b.header.Time = b.parent.Time() + uint64(r)
	// 		b.header.Difficulty = b.engine.CalcDifficulty(chainreader, b.header.Time, b.parent.Header())
	// 	}
	// }

	type ratioComparison struct {
		tdRatio float64
		penalty float64
	}
	gotRatioComparisons := []ratioComparison{}

	for i, c := range cases {

		if err := chain.Reset(); err != nil {
			t.Fatal(err)
		}
		easy, _ := GenerateChain(genesis.Config, genesisB, engine, db, c.easyLen, func(i int, b *BlockGen) {
			b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
			// poissonTime(b, c.easyOffset)
			b.OffsetTime(c.easyOffset - 10)
		})
		commonAncestor := easy[c.commonAncestorN-1]
		hard, _ := GenerateChain(genesis.Config, commonAncestor, engine, db, c.hardLen, func(i int, b *BlockGen) {
			b.SetNonce(types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64))))
			// poissonTime(b, c.hardOffset)
			b.OffsetTime(c.hardOffset - 10)
		})
		if _, err := chain.InsertChain(easy); err != nil {
			t.Fatal(err)
		}
		n, err := chain.InsertChain(hard)
		hardHead := chain.CurrentBlock().Hash() == hard[len(hard)-1].Hash()

		commons := plotter.XYs{}
		easys := plotter.XYs{}
		hards := plotter.XYs{}
		tdrs := plotter.XYs{}
		antigravities := plotter.XYs{}
		antigravities2 := plotter.XYs{}

		balance := plotter.XYs{}

		for i := 0; i < c.easyLen; i++ {
			td := chain.GetTd(easy[i].Hash(), easy[i].NumberU64())
			point := plotter.XY{X: float64(easy[i].NumberU64()), Y: float64(td.Uint64())}
			if i <= c.commonAncestorN {
				commons = append(commons, point)
			} else {
				easys = append(easys, point)
			}
		}
		// td ratios
		// for j := 0; j < c.hardLen; j++ {
		for j := 0; j < n; j++ {

			td := chain.GetTd(hard[j].Hash(), hard[j].NumberU64())
			if td != nil {
				point := plotter.XY{X: float64(hard[j].NumberU64()), Y: float64(td.Uint64())}
				hards = append(hards, point)
			}

			if commonAncestor.NumberU64() != uint64(c.commonAncestorN) {
				t.Fatalf("bad test common=%d easy=%d can=%d", commonAncestor.NumberU64(), c.easyLen, c.commonAncestorN)
			}

			ee := c.commonAncestorN + j
			easyHeader := easy[ee].Header()
			hardHeader := hard[j].Header()
			if easyHeader.Number.Uint64() != hardHeader.Number.Uint64() {
				t.Fatalf("bad test easyheader=%d hardheader=%d", easyHeader.Number.Uint64(), hardHeader.Number.Uint64())
			}

			/*
				HERE LIES THE RUB (IN MY GRAPHS).


			*/
			// y := chain.getTDRatio(commonAncestor.Header(), easyHeader, hardHeader) // <- unit x unit

			// y := chain.getTDRatio(commonAncestor.Header(), easy[c.easyLen-1].Header(), hardHeader)

			y := chain.getTDRatio(commonAncestor.Header(), chain.CurrentHeader(), hardHeader)

			if j == 0 {
				t.Logf("case=%d first.hard.tdr=%v", i, y)
			}

			ecbp := ecbp1100AGSinusoidalA(float64(hardHeader.Time - commonAncestor.Header().Time))

			if j == n-1 {
				gotRatioComparisons = append(gotRatioComparisons, ratioComparison{
					tdRatio: y, penalty: ecbp,
				})
			}

			// Exploring alternative penalty functions.
			ecbp2 := ecbp1100AGExpA(float64(hardHeader.Time - commonAncestor.Header().Time))
			// t.Log(y, ecbp, ecbp2)

			tdrs = append(tdrs, plotter.XY{X: float64(hard[j].NumberU64()), Y: y})
			antigravities = append(antigravities, plotter.XY{X: float64(hard[j].NumberU64()), Y: ecbp})
			antigravities2 = append(antigravities2, plotter.XY{X: float64(hard[j].NumberU64()), Y: ecbp2})

			balance = append(balance, plotter.XY{X: float64(hardHeader.Number.Uint64()), Y: y - ecbp})
		}
		scatterCommons, _ := plotter.NewScatter(commons)
		scatterEasys, _ := plotter.NewScatter(easys)
		scatterHards, _ := plotter.NewScatter(hards)

		scatterTDRs, _ := plotter.NewScatter(tdrs)
		scatterAntigravities, _ := plotter.NewScatter(antigravities)
		scatterAntigravities2, _ := plotter.NewScatter(antigravities2)
		balanceScatter, _ := plotter.NewScatter(balance)

		scatterCommons.Color = color.RGBA{R: 190, G: 197, B: 236, A: 255}
		scatterCommons.Shape = draw.CircleGlyph{}
		scatterCommons.Radius = 2
		scatterEasys.Color = color.RGBA{R: 152, G: 236, B: 161, A: 255} // green
		scatterEasys.Shape = draw.CircleGlyph{}
		scatterEasys.Radius = 2
		scatterHards.Color = color.RGBA{R: 236, G: 106, B: 94, A: 255}
		scatterHards.Shape = draw.CircleGlyph{}
		scatterHards.Radius = 2

		p, perr := plot.New()
		if perr != nil {
			log.Panic(perr)
		}
		p.Add(scatterCommons)
		p.Legend.Add("Commons", scatterCommons)
		p.Add(scatterEasys)
		p.Legend.Add("Easys", scatterEasys)
		p.Add(scatterHards)
		p.Legend.Add("Hards", scatterHards)
		p.Title.Text = fmt.Sprintf("TD easy=%d hard=%d", c.easyOffset, c.hardOffset)
		p.Save(1000, 600, fmt.Sprintf("plot-td-%d-%d-%d-%d-%d.png", c.easyLen, c.commonAncestorN, c.hardLen, c.easyOffset, c.hardOffset))

		p, perr = plot.New()
		if perr != nil {
			log.Panic(perr)
		}

		scatterTDRs.Color = color.RGBA{R: 236, G: 106, B: 94, A: 255} // red
		scatterTDRs.Radius = 3
		scatterTDRs.Shape = draw.PyramidGlyph{}
		p.Add(scatterTDRs)
		p.Legend.Add("TD Ratio", scatterTDRs)

		scatterAntigravities.Color = color.RGBA{R: 190, G: 197, B: 236, A: 255} // blue
		scatterAntigravities.Radius = 3
		scatterAntigravities.Shape = draw.PlusGlyph{}
		p.Add(scatterAntigravities)
		p.Legend.Add("(Anti)Gravity Penalty", scatterAntigravities)

		scatterAntigravities2.Color = color.RGBA{R: 152, G: 236, B: 161, A: 255} // green
		scatterAntigravities2.Radius = 3
		scatterAntigravities2.Shape = draw.PlusGlyph{}
		// p.Add(scatterAntigravities2)
		// p.Legend.Add("(Anti)Gravity Penalty (Alternate)", scatterAntigravities2)

		p.Title.Text = fmt.Sprintf("TD Ratio easy=%d hard=%d", c.easyOffset, c.hardOffset)
		p.Save(1000, 600, fmt.Sprintf("plot-td-ratio-%d-%d-%d-%d-%d.png", c.easyLen, c.commonAncestorN, c.hardLen, c.easyOffset, c.hardOffset))

		p, perr = plot.New()
		if perr != nil {
			log.Panic(perr)
		}
		p.Title.Text = fmt.Sprintf("TD Ratio - Antigravity Penalty easy=%d hard=%d", c.easyOffset, c.hardOffset)
		balanceScatter.Color = color.RGBA{R: 235, G: 92, B: 236, A: 255} // purple
		balanceScatter.Radius = 3
		balanceScatter.Shape = draw.PlusGlyph{}
		p.Add(balanceScatter)
		p.Legend.Add("TDR - Penalty", balanceScatter)
		p.Save(1000, 600, fmt.Sprintf("plot-td-ratio-diff-%d-%d-%d-%d-%d.png", c.easyLen, c.commonAncestorN, c.hardLen, c.easyOffset, c.hardOffset))

		if (err != nil && c.accepted) || (err == nil && !c.accepted) || (hardHead != c.hardGetsHead) {
			compared := gotRatioComparisons[i]
			t.Errorf(`case=%d [easy=%d hard=%d ca=%d eo=%d ho=%d] want.accepted=%v want.hardHead=%v got.hardHead=%v err=%v
got.tdr=%v got.pen=%v`,
				i,
				c.easyLen, c.hardLen, c.commonAncestorN, c.easyOffset, c.hardOffset,
				c.accepted, c.hardGetsHead, hardHead, err, compared.tdRatio, compared.penalty)
		}
	}

}
