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

var Blank = ""
var Manifest string
var UriPrefix string
var SubGroup string
var x264Profiles = map[string]string{"Baseline": "42", "Main": "4d", "High": "64"}
var AudioProfiles = map[string]string{"HE-AACv2": "mp4a.40.5", "LC": "mp4a.40.2", "mp3": "mp4a.40.34"}

type Format struct {
	FormatName string `json:"format_name"`
	BitRate    string `json:"bit_rate"`
}

type Stream struct {
	CodecType string  `json:"codec_type"`
	CodecName string  `json:"codec_name"`
	Profile   string  `json:"profile"`
	Level     float64 `json:"level"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
}

type Container struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}

type Stanza struct {
	Bandwidth  string
	Resolution string
	ACodec     string
	VCodec     string
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

//ensure urlprefix ends in a "/"
func fixPrefix(manifest string, uriprefix string) string {
	if uriprefix != Blank {
		if !(strings.HasSuffix(uriprefix, "/")) {
			if !(strings.HasPrefix(manifest, "/")) {
				uriprefix += "/"
			}
		}
	}
	return uriprefix
}

// create a subtitle stanza for use in the  master.m3u8
func mkSubStanza(manifest string, uriprefix string, subgroup string) string {
	if subgroup == Blank {
		subgroup = "WebVtt"
	}
	one := fmt.Sprintf("#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID=\"%s\",", subgroup)
	two := "NAME=\"English\",DEFAULT=YES,AUTOSELECT=YES,FORCED=NO,"
	three := fmt.Sprintf("LANGUAGE=\"en\",URI=\"%s%s\"\n", uriprefix, manifest)
	return one + two + three
}

// find the first segment in the m3u8 file
func findSegment(manifest string) string {
	file, err := os.Open(manifest)
	chk(err, "trouble reading manifest")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !(strings.HasPrefix(line, "#")) {
			segment := strings.Replace(manifest, path.Base(manifest), line, 1)

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

// determine audio codec for a stream
func setACodec(i Stream) string {
	if AudioProfiles[i.CodecName] != "" {
		return AudioProfiles[i.CodecName]

	}
	if AudioProfiles[i.Profile] != "" {
		return AudioProfiles[i.Profile]
	}
	return ""
}

// determine final codec value for stanza
func setStanzaCodec(st Stanza) string {
	if st.ACodec != "" {
		return fmt.Sprintf("\"%s\"", st.ACodec)
		if st.VCodec != "" {
			return fmt.Sprintf("\"%s,%s\"", st.VCodec, st.ACodec)
		}
	} else {
		if st.VCodec != "" {
			return fmt.Sprintf("\"%s\"", st.VCodec)
		}
	}
	return ""
}

//Generate full stanza for master.m3u8 file
func mkStanza(manifest string, segment string, subgroup string, uriprefix string) {
	var st Stanza
	var f Container
	jason := Probe(segment)
	err := json.Unmarshal(jason, &f)
	chk(err, "bad data while probing file")
	st.Bandwidth = f.Format.BitRate
	uriprefix = fixPrefix(manifest, uriprefix)
	st.ACodec = ""
	st.VCodec = ""
	codec := ""
	for _, i := range f.Streams {
		fmt.Println(i.CodecName)
		if i.CodecType == "subtitle" {
			substanza := mkSubStanza(manifest, uriprefix, subgroup)
			showStanza(substanza, Blank)
			return
		}
		if i.CodecType == "video" {
			st.Resolution = fmt.Sprintf("%vx%v", i.Width, i.Height)
			st.VCodec = fmt.Sprintf("avc1.%v00%x", x264Profiles[i.Profile], int(i.Level))
		}
		if i.CodecType == "audio" {
			st.ACodec = setACodec(i)
		}
	}
	codec = setStanzaCodec(st)
	m3u8Stanza := ""
	if st.VCodec != "" {
		m3u8Stanza = fmt.Sprintf("#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=%v,RESOLUTION=%s,CODECS=%s", st.Bandwidth, st.Resolution, codec)
		if subgroup != Blank {
			m3u8Stanza = fmt.Sprintf("%s,SUBTITLES=\"%s\"", m3u8Stanza, subgroup)
		}
	} else {
		m3u8Stanza = fmt.Sprintf("#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=%v,CODECS=%s", st.Bandwidth, codec)
	}
	mpath := fmt.Sprintf("%s%s\n", uriprefix, manifest)
	showStanza(m3u8Stanza, mpath)
}

// Command line flags
func mkFlags() {
	flag.StringVar(&Manifest, "i", Blank, "manifest file (required)")
	flag.StringVar(&SubGroup, "s", Blank, "add subtitle group i.e. SUBTITLES= (optional)")
	flag.StringVar(&UriPrefix, "u", Blank, "url prefix to add to index.m3u8 path (optional)")
	flag.Parse()
}

// find a segment, make a stanza
func do(manifest string, subgroup string, uriprefix string) {
	segment := findSegment(Manifest)
	mkStanza(manifest, segment, subgroup, uriprefix)
}

func main() {
	mkFlags()
	if Manifest != Blank {
		do(Manifest, SubGroup, UriPrefix)
	} else {
		flag.PrintDefaults()
	}
}
