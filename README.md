# Adaptive Bitrate Streaming using HLS and MPEG Dash in GO

## Create new project

- `mkdir -p go-video-streamer`
- `cd go-video-streamer`
- `go mod init go-video-streamer`

## Dependencies

- `go get -u github.com/gin-gonic/gin`
- `go get github.com/tkanos/gonfig`

## To build and run

### Build

- `go build .`

### Upgrade dependencies

- `go get -u`
- `go mod tidy`

### Run

- `./go-video-streamer`

## FFMPEG HLS Generator

`https://ottverse.com/hls-packaging-using-ffmpeg-live-vod/`
`https://ottverse.com/free-hls-m3u8-test-urls/`

```bash
ffmpeg -i 1280.mp4 -codec: copy -bsf:v h264_mp4toannexb -start_number 0 -hls_time 4 -hls_list_size 0 -f hls 1280.m3u8

ffmpeg -i caminandes_llamigos_1080p_hevc.mp4 \
  -filter_complex \
  "[0:v]split=3[v1][v2][v3]; \
  [v1]copy[v1out]; [v2]scale=w=1280:h=720[v2out]; [v3]scale=w=640:h=360[v3out]" \
  -map "[v1out]" -c:v:0 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:0 5M -maxrate:v:0 5M -minrate:v:0 5M -bufsize:v:0 10M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
  -map "[v2out]" -c:v:1 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:1 3M -maxrate:v:1 3M -minrate:v:1 3M -bufsize:v:1 3M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
  -map "[v3out]" -c:v:2 libx264 -x264-params "nal-hrd=cbr:force-cfr=1" -b:v:2 1M -maxrate:v:2 1M -minrate:v:2 1M -bufsize:v:2 1M -preset slow -g 48 -sc_threshold 0 -keyint_min 48 \
  -map a:0 -c:a:0 aac -b:a:0 96k -ac 2 \
  -map a:0 -c:a:1 aac -b:a:1 96k -ac 2 \
  -map a:0 -c:a:2 aac -b:a:2 48k -ac 2 \
  -f hls \
  -hls_time 10 \
  -hls_playlist_type vod \
  -hls_flags independent_segments \
  -hls_segment_type mpegts \
  -hls_segment_filename stream_%v/data%02d.ts \
  -master_pl_name master.m3u8 \
  -var_stream_map "v:0,a:0 v:1,a:1 v:2,a:2" stream_%v/stream.m3u8
```

## mp4dash DASH + HLS Generator

`https://ottverse.com/bento4-mp4dash-for-mpeg-dash-packaging/`
`https://github.com/axiomatic-systems/Bento4/issues/85`
`https://superuser.com/questions/908280/what-is-the-correct-way-to-fix-keyframes-in-ffmpeg-for-dash`

```bash
mp4fragment --fragment-duration 4000 caminandes_llamigos_1080p_hevc.mp4 frag_caminandes_llamigos_1080p_hevc.mp4
mp4info frag_caminandes_llamigos_1080p_hevc.mp4
mp4dash --mpd-name=master.mpd frag_caminandes_llamigos_1080p_hevc.mp4 --hls --hls-master-playlist-name=master.m3u8
ffmpeg -i caminandes_llamigos_1080p_hevc.mp4 -vf scale=1280:720 -preset slow -crf 18 caminandes_llamigos_720p_hevc.mp4

-1 to automatically compute the aspect ratio
ffmpeg -i caminandes_llamigos_1080p_hevc.mp4 -vf scale=640:-1 -preset slow -crf 18 caminandes_llamigos_360p.mp4
ffmpeg -i caminandes_llamigos_1080p_hevc.mp4 -vf scale=1280:720 -preset slow -crf 18 caminandes_llamigos_720p.mp4
ffmpeg -i caminandes_llamigos_360p.mp4 -vcodec libx265 -crf 28 caminandes_llamigos_360p_hevc.mp4

ffmpeg -i caminandes_llamigos_1080p_hevc_original.mp4 -vcodec libx265 -vf scale=1280:720 caminandes_llamigos_720p_hevc.mp4
ffmpeg -i caminandes_llamigos_1080p_hevc_original.mp4 -vcodec libx265 -vf scale=1920:1080 caminandes_llamigos_1080p_hevc.mp4
ffmpeg -i caminandes_llamigos_1080p_hevc_original.mp4 -vcodec libx265 -vf scale=640:360 caminandes_llamigos_360p_hevc.mp4

mp4fragment caminandes_llamigos_1080p_hevc.mp4 frag_caminandes_llamigos_1080p_hevc.mp4
mp4fragment caminandes_llamigos_720p_hevc.mp4 frag_caminandes_llamigos_720p_hevc.mp4
mp4fragment caminandes_llamigos_360p_hevc.mp4 frag_caminandes_llamigos_360p_hevc.mp4

mp4dash --mpd-name=master.mpd frag_caminandes_llamigos_1080p_hevc.mp4 frag_caminandes_llamigos_720p_hevc.mp4 frag_caminandes_llamigos_360p_hevc.mp4 --hls --hls-master-playlist-name=master.m3u8
```

## URLs to test streaming

### HLS

- `http://localhost:8080/browse/hls/master.m3u8`

### DASH

- `http://localhost:8080/browse/hls/stream.mpd`
