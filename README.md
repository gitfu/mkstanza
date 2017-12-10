# mkstanza
Generate variant stanza for master.m3u8 file 

## Install
* Requires Go and ffprobe
* git clone https://github.com/gitfu/mkstanza.git
* cd mkstanza
* go build mkstanza.go

## Usage
```go
 ./mkstanza 
  -i string
    	manifest file (required, an m3u8 file)
  -s string
    	add subtitle group i.e. SUBTITLES= (optional)
  -u string
    	url prefix to add to index.m3u8 path (optional)
```
## Examples

* (Audio Codec mp3 )
```sh 
./mkstanza -i mp3.m3u8
 
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=154074,CODECS="mp4a.40.34"
mp3.m3u8

```
* (Audio Codec: mp3  )
```go
./mkstanza -i mp3.m3u8 -u http://example.com

#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=154074,CODECS="mp4a.40.34"
http://example.com/mp3.m3u8

```
* ( Video Codec h264; profile High ; level 3.1 Audio Codec aac; profile HE-AACv2 )

```go 
./mkstanza  -i audio_and_video.m3u8

#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1994969,RESOLUTION=1280x720,CODECS="avc1.64001f,mp4a.40.5"
audio_and_video.m3u8

```
* ( Video Codec h264; profile High ; level 3 Audio Codec aac; profile HE-AACv2 )
```go
 ./mkstanza  -i audio_and_video.m3u8 - u http://example.com 

#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=866368,RESOLUTION=640x360,CODECS="avc1.64001e,mp4a.40.5"
http://example.com/audio_and_video.m3u8

```
* ( Video Codec h264; profile Main ; level 3.1  Audio Codec aac; profile LC )
```go
./mkstanza  -i audio_and_video.m3u8 - u http://example.com -s mySubGroup


#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1994969,RESOLUTION=1280x720,CODECS="avc1.4d001f,mp4a.40.2",SUBTITLES="mySubGroup"
http://example.com/audio_and_video.m3u8

```
* ( Subtitle Codec webvtt )

```go
./mkstanza -i index_vtt.m3u8


#EXT-XMEDIA:TYPE=SUBTITLES,GROUPID="WebVtt",NAME="Eng",DEFAULT=YES,AUTOSELECT=YES,FORCED=NO,LANGUAGE="en", URI="ndex_vtt.m3u8"

```


* ( Multiple input files, (-s) SubGroup= "fu", (-u) UriPrefix= "http://fu.zu")
```go
./mkstanza -i /home/leroy/manifesto/ctrl/720/index.m3u8 -i  /home/leroy/manifesto/ctrl/360/index.m3u8 -s "fu" -u http://fu.zu  -i /home/leroy/manifesto/99/360/index.m3u8   -i /home/leroy/manifesto/99/720/index.m3u8  -i /home/leroy/manifesto/99/432/index.m3u8 -i /home/leroy/manifesto/scte35/subs/index_vtt.m3u8  



#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=2191704,RESOLUTION=1280x720,CODECS="avc1.64001f,mp4a.40.5",SUBTITLES="fu"
http://fu.zu/home/leroy/manifesto/ctrl/720/index.m3u8


#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=966813,RESOLUTION=640x360,CODECS="avc1.64001e,mp4a.40.5",SUBTITLES="fu"
http://fu.zu/home/leroy/manifesto/ctrl/360/index.m3u8


#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=967302,RESOLUTION=640x360,CODECS="avc1.64001e,mp4a.40.5",SUBTITLES="fu"
http://fu.zu/home/leroy/manifesto/99/360/index.m3u8


#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=2189505,RESOLUTION=1280x720,CODECS="avc1.64001f,mp4a.40.5",SUBTITLES="fu"
http://fu.zu/home/leroy/manifesto/99/720/index.m3u8


#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1885970,RESOLUTION=768x432,CODECS="avc1.64001e,mp4a.40.5",SUBTITLES="fu"
http://fu.zu/home/leroy/manifesto/99/432/index.m3u8


#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID="fu",NAME="English",DEFAULT=YES,AUTOSELECT=YES,FORCED=NO,LANGUAGE="en",URI="http://fu.zu/home/leroy/manifesto/scte35/subs/index_vtt.m3u8"


```
