# mkstanza
Generate variant stanza for master.m3u8 file 

## Install
* Requires Go and ffprobe
* git clone https://github.com/gitfu/mkstanza.git
* cd mkstanza
* go build mkstanza.go


#### ./mkstanza -i mp3.m3u8
##### (Audio Codec mp3 )
```
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=154074,CODECS="mp4a.40.34"
mp3.m3u8
```

#### ./mkstanza -i mp3.m3u8 -u http://example.com
##### ***(Audio Codec mp3 )
```
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=154074,CODECS="mp4a.40.34"
http://example.com/mp3.m3u8
```
#### ./mkstanza  -i audio_and_video.m3u8
##### (Video Codec h264; profile High ; level 3.1 Audio Codec aac; profile HE-AACv2)

```
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1994969,RESOLUTION=1280x720,CODECS="avc1.64001f,mp4a.40.5"
audio_and_video.m3u8
```
#### ./mkstanza  -i audio_and_video.m3u8 - u http://example.com 
##### (Video Codec h264; profile High ; level 3 Audio Codec aac; profile HE-AACv2)


```
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=866368,RESOLUTION=640x360,CODECS="avc1.64001e,mp4a.40.5"
http://example.com/audio_and_video.m3u8
```

#### ./mkstanza  -i audio_and_video.m3u8 - u http://example.com -s mySubGroup
##### (Video Codec h264; profile Main ; level 3.1  Audio Codec aac; profile LC)

```
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1994969,RESOLUTION=1280x720,CODECS="avc1.4d001f,mp4a.40.2",SUBTITLES="mySubGroup"
http://example.com/audio_and_video.m3u8
```


#### ./mkstanza -i index_vtt.m3u8
```
#EXT-XMEDIA:TYPE=SUBTITLES,GROUPID="WebVtt",NAME="Eng",DEFAULT=YES,AUTOSELECT=YES,FORCED=NO,LANGUAGE="en", URI="ndex_vtt.m3u8"
```
