package main

/*

Author Gaurav Sablok
Universitat Potsdam
Date 2024-10-16

It takes all forms of the micrRNA predictions and all the target analyzer for the microRNA predictions
and prepares them for the deep learning.

*/

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

var (
	psRNAPred     string
	tapirPred     string
	fastPred      string
	evalue        float64
	psRNAfile     string
	tarHunter     string
	targetFinder  string
	tarFinderFile string
	upstream      int
	downstream    int
	psRobotFile   string
)

var rootCmd = &cobra.Command{
	Use:  "analyzePred",
	Long: "This analyzes the microRNA prediction and makes them ready for the deep learning approaches",
}

var psRNACmd = &cobra.Command{
	Use:  "psRNAanalyzer",
	Long: "Analyzes and prepares the psRNA target predictions for the deep learning",
	Run:  psRNAFunc,
}

var tapirCmd = &cobra.Command{
	Use:  "tapiranalyzer",
	Long: "Analyzes tapir target predictions for the deep learning",
	Run:  tapirFunc,
}

var psRNAMapCmd = &cobra.Command{
	Use:  "psRNAmapanalyze",
	Long: "Analyze the psRNA map alignment of the reads to the genome",
	Run:  psRNAMapFunc,
}

var tarHunterCmd = &cobra.Command{
	Use:  "tarHunter",
	Long: "analyze the tarHunter results for the miRNA predictions",
	Run:  tarFunc,
}

var tarFinderCmd = &cobra.Command{
	Use:  "targetFinder",
	Long: "analyzes the targetFinder results for the miRNA predictions",
	Run:  tarFinderFunc,
}

var psRobotCmd = &cobra.Command{
	Use:  "psRobot",
	Long: "analyzes the results from the psRobot microRNA predictions",
	Run:  psRobotFunc,
}

func init() {
	psRNACmd.Flags().
		StringVarP(&psRNAPred, "psRNAPred", "p", "psRNA microRNA predictions", "psRNA predictions")
	psRNACmd.Flags().
		StringVarP(&fastPred, "fastapred", "f", "fasta file for the predictions", "fasta predict")
	psRNACmd.Flags().
		Float64Var(&evalue, "expectation value", 0.5, "expectation value")
	psRNACmd.Flags().
		IntVarP(&upstream, "upstream", "U", 10, "upstream of the miRNA predictions")
	psRNACmd.Flags().
		IntVarP(&downstream, "downstream", "D", 10, "downstream of the miRNA predictions")
	tapirCmd.Flags().
		StringVarP(&psRNAPred, "tapir", "p", "psRNA microRNA predictions", "psRNA predictions")
	tapirCmd.Flags().
		StringVarP(&fastPred, "fastapred", "f", "fasta file for the predictions", "fasta predict")
	tapirCmd.Flags().
		IntVarP(&upstream, "upstream", "U", 10, "upstream of the miRNA predictions")
	tapirCmd.Flags().
		IntVarP(&downstream, "downstream", "D", 10, "downstream of the miRNA predictions")
	psRNAMapCmd.Flags().
		StringVarP(&psRNAfile, "psRNAmapfile", "P", "RNA mapping file", "map reads to the genome file")
	psRNAMapCmd.Flags().
		StringVarP(&fastPred, "fastapred", "f", "fasta file for the predictions", "fasta predict")
	psRNAMapCmd.Flags().
		IntVarP(&upstream, "upstream", "U", 10, "upstream of the miRNA predictions")
	psRNAMapCmd.Flags().
		IntVarP(&downstream, "downstream", "D", 10, "downstream of the miRNA predictions")
	tarHunterCmd.Flags().
		StringVarP(&tarHunter, "tarhunterfile", "T", "tarHunter predictions", "tarhunter analysis")
	tarHunterCmd.Flags().
		StringVarP(&fastPred, "fastapred", "f", "fasta file for the predictions", "fasta predict")
	tarHunterCmd.Flags().
		IntVarP(&upstream, "upstream", "U", 10, "upstream of the miRNA predictions")
	tarHunterCmd.Flags().
		IntVarP(&downstream, "downstream", "D", 10, "downstream of the miRNA predictions")
	tarFinderCmd.Flags().
		StringVarP(&tarFinderFile, "targetFinderfile", "T", "targetFinder predictions", "targetFinder analysis")
	tarFinderCmd.Flags().
		StringVarP(&fastPred, "fastapred", "f", "fasta file for the predictions", "fasta predict")
	tarFinderCmd.Flags().
		IntVarP(&upstream, "upstream", "U", 10, "upstream of the miRNA predictions")
	tarFinderCmd.Flags().
		IntVarP(&downstream, "downstream", "D", 10, "downstream of the miRNA predictions")
	psRobotCmd.Flags().
		StringVarP(&psRobotFile, "psRobotfile", "R", "psRobot predictions", " psRobot analysis")
	psRobotCmd.Flags().
		StringVarP(&fastPred, "fastapred", "f", "fasta file for the predictions", "fasta predict")
	psRobotCmd.Flags().
		IntVarP(&upstream, "upstream", "U", 10, "upstream of the miRNA predictions")
	psRobotCmd.Flags().
		IntVarP(&downstream, "downstream", "D", 10, "downstream of the miRNA predictions")

	rootCmd.AddCommand(psRNACmd)
	rootCmd.AddCommand(tapirCmd)
	rootCmd.AddCommand(psRNAMapCmd)
	rootCmd.AddCommand(tarHunterCmd)
	rootCmd.AddCommand(tarFinderCmd)
	rootCmd.AddCommand(psRobotCmd)
}

func psRNAFunc(cmd *cobra.Command, args []string) {
	type psRNAStruct struct {
		miRNA         string
		target        string
		evalue        float64
		targetStart   int
		targetEnd     int
		miRNAAligned  string
		targetAligned string
	}

	type psRNAStructFiltered struct {
		miRNA         string
		target        string
		evalue        float64
		targetStart   int
		targetEnd     int
		miRNAAligned  string
		targetAligned string
	}

	type fastPredID struct {
		id string
	}

	type fastPredSeq struct {
		seq string
	}

	fastID := []fastPredID{}
	fastSeq := []fastPredSeq{}

	fastaOpen, err := os.Open(fastPred)
	if err != nil {
		log.Fatal(err)
	}
	fastaRead := bufio.NewScanner(fastaOpen)
	for fastaRead.Scan() {
		line := fastaRead.Text()
		if strings.HasPrefix(string(line), ">") {
			fastID = append(fastID, fastPredID{
				id: strings.ReplaceAll(string(line), ">", ""),
			})
		}
		if !strings.HasPrefix(string(line), ">") {
			fastSeq = append(fastSeq, fastPredSeq{
				seq: string(line),
			})
		}
	}

	storemiRNA := []psRNAStruct{}
	filteredmiRNA := []psRNAStructFiltered{}

	fOpen, err := os.Open(psRNAPred)
	if err != nil {
		log.Fatal(err)
	}
	fRead := bufio.NewScanner(fOpen)
	for fRead.Scan() {
		line := fRead.Text()
		if strings.HasPrefix(string(line), "#") {
			continue
		}
		if !strings.HasPrefix(string(line), "#") {
			evalueMid, _ := strconv.ParseFloat(strings.Split(string(line), " ")[2], 64)
			targetStartMid, _ := strconv.Atoi(strings.Split(string(line), " ")[6])
			targetEndMid, _ := strconv.Atoi(strings.Split(string(line), " ")[7])
			storemiRNA = append(storemiRNA, psRNAStruct{
				miRNA:         strings.Split(string(line), " ")[0],
				target:        strings.Split(string(line), " ")[1],
				evalue:        evalueMid,
				targetStart:   targetStartMid,
				targetEnd:     targetEndMid,
				miRNAAligned:  strings.Split(string(line), " ")[8],
				targetAligned: strings.Split(string(line), " ")[9],
			})
		}
	}

	for i := range storemiRNA {
		if storemiRNA[i].evalue <= evalue {
			filteredmiRNA = append(filteredmiRNA, psRNAStructFiltered{
				miRNA:         storemiRNA[i].miRNA,
				target:        storemiRNA[i].target,
				evalue:        storemiRNA[i].evalue,
				targetStart:   storemiRNA[i].targetStart,
				targetEnd:     storemiRNA[i].targetEnd,
				miRNAAligned:  storemiRNA[i].miRNAAligned,
				targetAligned: storemiRNA[i].targetAligned,
			})
		}
	}

	type extractSeq struct {
		target      string
		targetSeq   string
		targetcomp  string
		targetStart int
		targetEnd   int
		upstream    string
		downstream  string
	}

	targetExtract := []extractSeq{}

	for i := range filteredmiRNA {
		for j := range fastID {
			if filteredmiRNA[i].target == fastID[j].id {
				targetExtract = append(targetExtract, extractSeq{
					target:      filteredmiRNA[i].target,
					targetSeq:   fastSeq[j].seq[filteredmiRNA[i].targetStart:filteredmiRNA[i].targetEnd],
					targetStart: filteredmiRNA[i].targetStart,
					targetEnd:   filteredmiRNA[i].targetEnd,
					targetcomp:  fastSeq[j].seq,
					upstream:    fastSeq[j].seq[filteredmiRNA[i].targetStart-upstream : filteredmiRNA[i].targetEnd],
					downstream:  fastSeq[j].seq[filteredmiRNA[i].targetEnd : filteredmiRNA[i].targetEnd+downstream],
				})
			}
		}
	}

	psRNAneural, err := os.Create("psRNANeural.fasta")
	if err != nil {
		log.Fatal(err)
	}
	defer psRNAneural.Close()
	for i := range targetExtract {
		psRNAneural.WriteString(
			targetExtract[i].target + "\t" + ">" + targetExtract[i].targetSeq + "\t" + targetExtract[i].targetcomp + "\t" + targetExtract[i].upstream + "\t" + targetExtract[i].downstream,
		)
	}
}

func tapirFunc(cmd *cobra.Command, args []string) {
	target := []string{}
	miRNA := []string{}
	score := []int{}
	mfe := []int{}
	start := []int{}
	end := []int{}

	tapirOpen, err := os.Open(tapirPred)
	if err != nil {
		log.Fatal(err)
	}
	tapirRead := bufio.NewScanner(tapirOpen)
	for tapirRead.Scan() {
		line := tapirRead.Text()
		if strings.HasPrefix(string(line), "target") {
			target = append(target, strings.Split(string(line), " ")[1])
		}
		if strings.HasPrefix(string(line), "miRNA") {
			miRNA = append(miRNA, strings.Split(strings.Split(string(line), ":")[3], "=")[2])
		}
		if strings.HasPrefix(string(line), "score") {
			scoreInter, _ := strconv.Atoi(strings.Split(string(line), " ")[1])
			score = append(score, scoreInter)
		}
		if strings.HasPrefix(string(line), "mfe") {
			mfeInter, _ := strconv.Atoi(strings.Split(string(line), " ")[1])
			mfe = append(mfe, mfeInter)
		}
		if strings.HasPrefix(string(line), "start") {
			startInter, _ := strconv.Atoi(strings.Split(string(line), " ")[1])
			start = append(start, startInter)
		}
		if strings.HasPrefix(string(line), "target_5'") {
			end = append(end, len(strings.Split(string(line), " ")[1]))
		}
	}

	type tapirExtract struct {
		tapid      string
		tapseq     string
		tapcomp    string
		upstream   string
		downstream string
	}

	type fastPredID struct {
		id string
	}

	type fastPredSeq struct {
		seq string
	}

	fastID := []fastPredID{}
	fastSeq := []fastPredSeq{}

	fastaOpen, err := os.Open(fastPred)
	if err != nil {
		log.Fatal(err)
	}
	fastaRead := bufio.NewScanner(fastaOpen)
	for fastaRead.Scan() {
		line := fastaRead.Text()
		if strings.HasPrefix(string(line), ">") {
			fastID = append(fastID, fastPredID{
				id: strings.ReplaceAll(string(line), ">", ""),
			})
		}
		if !strings.HasPrefix(string(line), ">") {
			fastSeq = append(fastSeq, fastPredSeq{
				seq: string(line),
			})
		}
	}

	tapSeq := []tapirExtract{}

	for i := range target {
		for j := range fastID {
			if target[i] == fastID[j].id {
				tapSeq = append(tapSeq, tapirExtract{
					tapid:      target[i],
					tapseq:     fastSeq[j].seq[start[i]:end[i]],
					tapcomp:    fastSeq[j].seq,
					upstream:   fastSeq[j].seq[start[i]-upstream : start[i]],
					downstream: fastSeq[j].seq[end[i] : end[i]+downstream],
				})
			}
		}
	}

	tapwrite, err := os.Create("tapirneural.fasta")
	if err != nil {
		log.Fatal(err)
	}
	defer tapwrite.Close()
	for i := range tapSeq {
		tapwrite.WriteString(
			tapSeq[i].tapid + "\t" + tapSeq[i].tapseq + "\t" + tapSeq[i].tapcomp + "\t" + tapSeq[i].upstream + "\t" + tapSeq[i].downstream,
		)
	}
}

func psRNAMapFunc(cmd *cobra.Command, args []string) {
	type psRNAfunc struct {
		id     string
		ref    string
		strand string
		start  int
		stop   int
		read   string
	}

	readMap := []psRNAfunc{}
	fOpen, err := os.Open(psRNAfile)
	if err != nil {
		log.Fatal(err)
	}
	fRead := bufio.NewScanner(fOpen)
	for fRead.Scan() {
		line := fRead.Text()
		startStore, _ := strconv.Atoi(strings.Split(string(line), " ")[3])
		endStore, _ := strconv.Atoi(strings.Split(string(line), " ")[4])
		readMap = append(readMap, psRNAfunc{
			id:     strings.Split(string(line), " ")[0],
			ref:    strings.Split(string(line), " ")[1],
			strand: strings.Split(string(line), " ")[2],
			start:  startStore,
			stop:   endStore,
			read:   strings.Split(string(line), " ")[5],
		})
	}

	type fastPredID struct {
		id string
	}

	type fastPredSeq struct {
		seq string
	}

	fastID := []fastPredID{}
	fastSeq := []fastPredSeq{}

	fastaOpen, err := os.Open(fastPred)
	if err != nil {
		log.Fatal(err)
	}
	fastaRead := bufio.NewScanner(fastaOpen)
	for fastaRead.Scan() {
		line := fastaRead.Text()
		if strings.HasPrefix(string(line), ">") {
			fastID = append(fastID, fastPredID{
				id: strings.ReplaceAll(string(line), ">", ""),
			})
		}
		if !strings.HasPrefix(string(line), ">") {
			fastSeq = append(fastSeq, fastPredSeq{
				seq: string(line),
			})
		}
	}

	type readMapSeq struct {
		read       string
		start      int
		stop       int
		refseq     string
		refcomp    string
		upstream   string
		downstream string
	}

	readMapSeqStore := []readMapSeq{}

	for i := range readMap {
		for j := range fastID {
			if readMap[i].id == fastID[i].id {
				readMapSeqStore = append(readMapSeqStore, readMapSeq{
					read:       readMap[i].read,
					start:      readMap[i].start,
					stop:       readMap[i].stop,
					refseq:     fastSeq[j].seq[readMap[i].start:readMap[i].stop],
					refcomp:    fastSeq[j].seq,
					upstream:   fastSeq[j].seq[readMap[i].start-upstream : readMap[i].start],
					downstream: fastSeq[j].seq[readMap[i].stop : readMap[i].stop+downstream],
				})
			}
		}
	}

	readMapWrite, err := os.Create("psRNAMap.fasta")
	if err != nil {
		log.Fatal(err)
	}
	for i := range readMapSeqStore {
		readMapWrite.WriteString(
			readMapSeqStore[i].read + "\t" + readMapSeqStore[i].refseq + "\t" + readMapSeqStore[i].refseq + "\t" + readMapSeqStore[i].refcomp + "\t" + readMapSeqStore[i].upstream + "\t" + readMapSeqStore[i].downstream,
		)
	}
}

func tarFunc(cmd *cobra.Command, args []string) {
	type tarAnalyze struct {
		targetID  string
		targetSeq string
		miRID     string
		miRSeq    string
		start     int
		end       int
	}

	tarStore := []tarAnalyze{}

	fOpen, err := os.Open(tarHunter)
	if err != nil {
		log.Fatal(err)
	}
	fRead := bufio.NewScanner(fOpen)
	for fRead.Scan() {
		line := fRead.Text()
		startAppend, _ := strconv.Atoi(strings.Split(string(line), " ")[9])
		endAppend, _ := strconv.Atoi(strings.Split(string(line), " ")[10])
		tarStore = append(tarStore, tarAnalyze{
			targetID:  strings.Split(string(line), " ")[0],
			targetSeq: strings.Split(string(line), " ")[1],
			miRID:     strings.Split(string(line), " ")[2],
			miRSeq:    strings.Split(string(line), " ")[3],
			start:     startAppend,
			end:       endAppend,
		})
	}

	type fastPredID struct {
		id string
	}

	type fastPredSeq struct {
		seq string
	}

	fastID := []fastPredID{}
	fastSeq := []fastPredSeq{}

	fastaOpen, err := os.Open(fastPred)
	if err != nil {
		log.Fatal(err)
	}
	fastaRead := bufio.NewScanner(fastaOpen)
	for fastaRead.Scan() {
		line := fastaRead.Text()
		if strings.HasPrefix(string(line), ">") {
			fastID = append(fastID, fastPredID{
				id: strings.ReplaceAll(string(line), ">", ""),
			})
		}
		if !strings.HasPrefix(string(line), ">") {
			fastSeq = append(fastSeq, fastPredSeq{
				seq: string(line),
			})
		}
	}

	type tarSeqStore struct {
		targetID         string
		targetSeq        string
		miRNAID          string
		miRANSeq         string
		targetSeqComp    string
		targetUpstream   string
		targetdownStream string
	}

	tarStoreCapture := []tarSeqStore{}

	for i := range tarStore {
		for j := range fastID {
			if tarStore[i].targetSeq == fastID[j].id {
				tarStoreCapture = append(tarStoreCapture, tarSeqStore{
					targetID:         tarStore[i].targetID,
					targetSeq:        tarStore[i].targetSeq,
					targetSeqComp:    fastSeq[j].seq,
					targetUpstream:   fastSeq[j].seq[tarStore[i].start : tarStore[i].start-upstream],
					targetdownStream: fastSeq[j].seq[tarStore[i].end : tarStore[i].end+downstream],
				})
			}
		}
	}

	tarHunterWrite, err := os.Create("tarHunter.fasta")
	if err != nil {
		log.Fatal(err)
	}
	for i := range tarStoreCapture {
		tarHunterWrite.WriteString(
			tarStoreCapture[i].targetID + "\t" + tarStoreCapture[i].targetSeq + "\t" + tarStoreCapture[i].targetSeq + tarStoreCapture[i].targetSeqComp + "\t" + tarStoreCapture[i].targetUpstream + "\t" + tarStoreCapture[i].targetdownStream,
		)
	}
}

func tarFinderFunc(cmd *cobra.Command, args []string) {
	type tarFinderCap struct {
		miRNA    string
		sequence string
		start    int
		end      int
		mfe      float64
	}

	tarFinderAdd := []tarFinderCap{}

	fOpen, err := os.Open(tarFinderFile)
	if err != nil {
		log.Fatal(err)
	}
	fRead := bufio.NewScanner(fOpen)
	for fRead.Scan() {
		line := fRead.Text()
		startConvert, _ := strconv.Atoi(strings.Split(string(line), " ")[15])
		endConvert, _ := strconv.Atoi(strings.Split(string(line), " ")[16])
		mfeConvert, _ := strconv.ParseFloat(strings.Split(string(line), " ")[18], 64)
		tarFinderAdd = append(tarFinderAdd, tarFinderCap{
			miRNA:    strings.Split(string(line), " ")[0],
			sequence: strings.Split(string(line), " ")[11],
			start:    startConvert,
			end:      endConvert,
			mfe:      mfeConvert,
		})
	}

	type fastPredID struct {
		id string
	}

	type fastPredSeq struct {
		seq string
	}

	fastID := []fastPredID{}
	fastSeq := []fastPredSeq{}

	fastaOpen, err := os.Open(fastPred)
	if err != nil {
		log.Fatal(err)
	}
	fastaRead := bufio.NewScanner(fastaOpen)
	for fastaRead.Scan() {
		line := fastaRead.Text()
		if strings.HasPrefix(string(line), ">") {
			fastID = append(fastID, fastPredID{
				id: strings.ReplaceAll(string(line), ">", ""),
			})
		}
		if !strings.HasPrefix(string(line), ">") {
			fastSeq = append(fastSeq, fastPredSeq{
				seq: string(line),
			})
		}
	}

	type tarFinderSeqStore struct {
		miRNA      string
		sequence   string
		compseq    string
		upstream   string
		downstream string
	}

	tarFinderSeqAdd := []tarFinderSeqStore{}
	for i := range tarFinderAdd {
		for j := range fastID {
			if tarFinderAdd[i].sequence == fastID[j].id {
				tarFinderSeqAdd = append(tarFinderSeqAdd, tarFinderSeqStore{
					miRNA:      tarFinderAdd[i].miRNA,
					sequence:   tarFinderAdd[i].sequence,
					compseq:    fastSeq[j].seq,
					upstream:   fastSeq[j].seq[tarFinderAdd[i].start-upstream : tarFinderAdd[i].start],
					downstream: fastSeq[j].seq[tarFinderAdd[i].end : tarFinderAdd[i].end+downstream],
				})
			}
		}
	}

	tarFinderWrite, err := os.Create("tarFinder.fasta")
	if err != nil {
		log.Fatal(err)
	}
	defer tarFinderWrite.Close()
	for i := range tarFinderSeqAdd {
		tarFinderWrite.WriteString(
			tarFinderSeqAdd[i].miRNA + "\t" + tarFinderSeqAdd[i].sequence + "\t" + tarFinderSeqAdd[i].compseq + "\t" + tarFinderSeqAdd[i].upstream + "\t" + tarFinderSeqAdd[i].downstream,
		)
	}
}

func psRobotFunc(cmd *cobra.Command, args []string) {
	type psRobotCapture struct {
		query        string
		score        float64
		queryStart   int
		queryEnd     int
		subjectStart int
		subjectEnd   int
	}

	psRobotC := []psRobotCapture{}
	sortedID := []string{}
	sortedScore := []float64{}

	sortedQueryStart := []int{}
	sortedQueryEnd := []int{}

	sortedSubjectStart := []int{}
	sortedSubjectEnd := []int{}

	fOpen, err := os.Open(psRobotFile)
	if err != nil {
		log.Fatal(err)
	}
	fRead := bufio.NewScanner(fOpen)
	for fRead.Scan() {
		line := fRead.Text()
		if strings.HasPrefix(string(line), ">") {
			capStore := strings.Split(string(line), "\t")[0]
			capScore := strings.Split(strings.Split(string(line), "\t")[1], ":")[1]
			scoreFloat, _ := strconv.ParseFloat(capScore, 64)
			sortedID = append(sortedID, strings.ReplaceAll(capStore, ">", ""))
			sortedScore = append(sortedScore, scoreFloat)
		}
		if strings.HasPrefix(string(line), "Query") {
			queryStart, _ := strconv.Atoi(
				strings.Split(strings.Split(string(line), "\t")[1], "\t")[0],
			)
			queryEnd, _ := strconv.Atoi(
				strings.Split(strings.Split(string(line), "\t")[1], "\t")[2],
			)
			sortedQueryStart = append(sortedQueryStart, queryStart)
			sortedQueryEnd = append(sortedQueryEnd, queryEnd)
		}
		if strings.HasPrefix(string(line), "Sbjct") {
			sbjctStart, _ := strconv.Atoi(
				strings.Split(strings.Split(string(line), "\t")[1], "\t")[0],
			)
			sbjctEnd, _ := strconv.Atoi(
				strings.Split(strings.Split(string(line), "\t")[1], "\t")[2],
			)
			sortedSubjectStart = append(sortedSubjectStart, sbjctStart)
			sortedSubjectEnd = append(sortedSubjectEnd, sbjctEnd)
		}
	}

	for i := range sortedID {
		psRobotC = append(psRobotC, psRobotCapture{
			query:        sortedID[i],
			score:        sortedScore[i],
			queryStart:   sortedQueryStart[i],
			queryEnd:     sortedQueryEnd[i],
			subjectStart: sortedSubjectStart[i],
			subjectEnd:   sortedSubjectEnd[i],
		})
	}

	type fastPredID struct {
		id string
	}

	type fastPredSeq struct {
		seq string
	}

	fastID := []fastPredID{}
	fastSeq := []fastPredSeq{}

	fastaOpen, err := os.Open(fastPred)
	if err != nil {
		log.Fatal(err)
	}
	fastaRead := bufio.NewScanner(fastaOpen)
	for fastaRead.Scan() {
		line := fastaRead.Text()
		if strings.HasPrefix(string(line), ">") {
			fastID = append(fastID, fastPredID{
				id: strings.ReplaceAll(string(line), ">", ""),
			})
		}
		if !strings.HasPrefix(string(line), ">") {
			fastSeq = append(fastSeq, fastPredSeq{
				seq: string(line),
			})
		}
	}

	type psRobotSeq struct {
		query         string
		score         float64
		subjectString string
		subjectComp   string
		upstream      string
		downstream    string
	}

	finalpsRobotSeq := []psRobotSeq{}

	for i := range psRobotC {
		for j := range fastID {
			if fastID[j].id == psRobotC[i].query {
				finalpsRobotSeq = append(finalpsRobotSeq, psRobotSeq{
					query:         psRobotC[i].query,
					score:         psRobotC[i].score,
					subjectString: fastSeq[j].seq[psRobotC[i].subjectStart:psRobotC[i].subjectEnd],
					subjectComp:   fastSeq[j].seq,
					upstream:      fastSeq[j].seq[psRobotC[i].subjectStart-upstream : psRobotC[i].subjectStart],
					downstream:    fastSeq[j].seq[psRobotC[i].subjectEnd : psRobotC[i].subjectEnd+downstream],
				})
			}
		}
	}
}
