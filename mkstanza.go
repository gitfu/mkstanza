package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

var SubGroup string
var UriPrefix string
var Blank = ""
//Maps for codec identifier strings
var x264Profiles = map[string]string{"Baseline": "42E0", "Main": "4d40", "High": "6400"}

var AudioProfiles = map[string]string{"HE-AAC": "mp4a.40.5","HE-AACv2": "mp4a.40.28",
				      "LC": "mp4a.40.2", "mp3": "mp4a.40.34"}

// To handle multiple inputs 
type infiles []string

// the flag.Value interface
func (nf *infiles) String() string {
	return fmt.Sprintf("%s", *nf)
}

// The second method is Set
func (nf *infiles) Set(value string) error {
	//fmt.Printf("%s\n", value)
	*nf = append(*nf, value)
	return nil
}

var nf infiles

// Struct for media format
type Format struct {
	FormatName string `json:"format_name"`
	BitRate    string `json:"bit_rate"`
}

// struct for individual streams
type Stream struct {
	CodecType string  `json:"codec_type"`
	CodecName string  `json:"codec_name"`
	Profile   string  `json:"profile"`
	Level     float64 `json:"level"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
}
// struct for media container
type Container struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}

// struct for variant stanza
type Stanza struct {
	Manifest   string
	UriPrefix  string
	SubGroup   string
	Bandwidth  string
	Resolution string
	ACodec     string
	VCodec     string
	Segment    string
}
// generate video codec identifier string
func (st *Stanza) SetVCodec(i Stream) {
	st.Resolution = fmt.Sprintf("%vx%v", i.Width, i.Height)
	if i.CodecName == "h264" {
		if x264Profiles[i.Profile] != Blank {
			st.VCodec = fmt.Sprintf("avc1.%v%x", x264Profiles[i.Profile], int(i.Level))
		}
	} else {
		unSupCodec(i.CodecName)
	}

}

// determine audio codec for a stream
func (st *Stanza) SetACodec(i Stream) {
	if AudioProfiles[i.CodecName] != Blank {
		st.ACodec = AudioProfiles[i.CodecName]
		return
	}
	if AudioProfiles[i.Profile] != Blank {
		st.ACodec = AudioProfiles[i.Profile]
		return
	}
	if st.ACodec == Blank {
		unSupCodec(i.CodecName)
	}
}

//ensure urlprefix ends in a "/"
func (st *Stanza) FixPrefix() {
	if st.UriPrefix != Blank {
		if !(strings.HasSuffix(st.UriPrefix, "/")) {
			if !(strings.HasPrefix(st.Manifest, "/")) {
				st.UriPrefix += "/"
			}
		}
	}

}

// create a subtitle stanza for use in the  master.m3u8
func (st *Stanza) mkSubStanza() string {
	if st.SubGroup == Blank {
		st.SubGroup = "WebVtt"
	}
	one := fmt.Sprintf("#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID=\"%s\",", st.SubGroup)
	two := "NAME=\"English\",DEFAULT=YES,AUTOSELECT=YES,FORCED=NO,"
	three := fmt.Sprintf("LANGUAGE=\"en\",URI=\"%s%s\"\n", st.UriPrefix, st.Manifest)
	return one + two + three
}

// determine final codec value for stanza
func (st *Stanza) CodecString() string {
	if st.VCodec != Blank && st.ACodec != Blank {
		return fmt.Sprintf("\"%s,%s\"", st.VCodec, st.ACodec)
	}
	if st.ACodec != Blank {
		return fmt.Sprintf("\"%s\"", st.ACodec)
	}
	if st.VCodec != Blank {
		return fmt.Sprintf("\"%s\"", st.VCodec)
	}
	return Blank
}
// handle unsupported codecs 
func unSupCodec(codecName string) {
	fmt.Printf("the codec %s is not currently supported.\n\n", codecName)
	syscall.Exit(-1)
}

// Generic catchall error checking
func chk(err error, mesg string) {
	if err != nil {
		fmt.Printf("%s\n", mesg)
		syscall.Exit(-1)
	}
}

// ffprobe a segment from the m3u8 file
func Probe(segment string) []byte {
	one := "ffprobe -hide_banner  -show_entries format=bit_rate -show_entries "
	two := "stream=codec_type,codec_name,height,width,profile,level -of json -i "
	cmd := fmt.Sprintf("%s%s%s", one, two, segment)
	parts := strings.Fields(cmd)
	data, err := exec.Command(parts[0], parts[1:]...).Output()
	chk(err, fmt.Sprintf("Error running \n %s \n %v", cmd, string(data)))
	return data
}

// find the first segment in the m3u8 file
func findSegment(manifest string) string {
	file, err := os.Open(manifest)
	chk(err, fmt.Sprintf("trouble reading %s", manifest))
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !(strings.HasPrefix(line, "#")) {
			// manifest="/hls/720/index.m3u8" path.Base(manifest)="index.m3u8", line="index0.ts"
			segment := strings.Replace(manifest, path.Base(manifest), line, 1)
			// segment ="/hls/720/index0.ts"
			return segment
		}
	}
	return Blank
}

// print m3u8  variant entry
func showStanza(stanza string, mpath string) {
	fmt.Println("")
	fmt.Println(stanza)
	fmt.Println(mpath)
}

//Generate full stanza for master.m3u8 file
func mkStanza(st Stanza) {
	var f Container
	jason := Probe(st.Segment)
	err := json.Unmarshal(jason, &f)
	chk(err, "bad data while probing file")
	st.Bandwidth = f.Format.BitRate
	codec := Blank
	for _, i := range f.Streams {
		//fmt.Printf("name %s; profile %s ; level %v\n",i.CodecName,i.Profile,i.Level)
		if i.CodecType == "subtitle" {
			substanza := st.mkSubStanza()
			showStanza(substanza, Blank)
			return
		}
		if i.CodecType == "video" {
			st.SetVCodec(i)
		}
		if i.CodecType == "audio" {
			st.SetACodec(i)
		}
	}
	codec = st.CodecString()
	m3u8Stanza := Blank
	if st.VCodec != Blank {
		m3u8Stanza = fmt.Sprintf("#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=%v,RESOLUTION=%s,CODECS=%s", st.Bandwidth, st.Resolution, codec)
		if st.SubGroup != Blank {
			m3u8Stanza = fmt.Sprintf("%s,SUBTITLES=\"%s\"", m3u8Stanza, st.SubGroup)
		}
	} else {
		m3u8Stanza = fmt.Sprintf("#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=%v,CODECS=%s", st.Bandwidth, codec)
	}
	mpath := fmt.Sprintf("%s%s\n", st.UriPrefix, st.Manifest)
	showStanza(m3u8Stanza, mpath)
}

// Makes it easy to call without command line flags/vars
func do(st Stanza) {
	if st.UriPrefix != Blank {
		st.FixPrefix()
	}
	st.Segment = findSegment(st.Manifest)
	mkStanza(st)
}

func main() {

	// Command line flags
	flag.Var(&nf, "i", " one or more input files(-i input1.m3u8 -i input2.m3u8 -i sub1.m3u8)")
	flag.StringVar(&SubGroup, "s", Blank, "add subtitle group i.e. SUBTITLES= (optional)")
	flag.StringVar(&UriPrefix, "u", Blank, "url prefix to add to index.m3u8 path (optional)")
	flag.Parse()
	if len(nf) == 0 {
		flag.PrintDefaults()
	} else {
		for i := 0; i < len(nf); i++ {
			var st Stanza
			st.SubGroup = SubGroup
			st.UriPrefix = UriPrefix
			st.Manifest = nf[i]
			if st.Manifest != Blank {
				do(st)
			}
		}
	}
}
