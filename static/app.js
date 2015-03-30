(function($) {
    
    $(function() {
        var url = location.pathname;
        var tracks = [];

        for (i=0; i<audio.length; i++) {
            var a = audio[i];
            var t = {
                title: a.Title
            };
            var u = "/_file" + url + "/" + a.File;
            if (a.Type == "mp3") {
                t.mp3 = u;
            }
            if (a.Type == "m4a") {
                t.m4a = u;
            }
            if (a.Type == "oga" || a.Type == "ogg") {
                t.oga = u;
            }
            tracks.push(t);
        }

        console.log(tracks);

        new jPlayerPlaylist({
            jPlayer: "#player",
            cssSelectorAncestor: "#playlist"
        }, tracks, {
            swfPath: "../../dist/jplayer",
            supplied: "oga, mp3, m4a",
            wmode: "window",
            useStateClassSkin: true,
            autoBlur: false,
            smoothPlayBar: true,
            keyEnabled: true
        });
    });

    // finally, initialize the player..
    $('#player').jPlayer({
     solution: 'html',
     supplied: 'oga, mp3, m4a',
     preload: 'metadata',
     volume: 0.8,
     muted: false,
     backgroundColor: '#000000',
     cssSelectorAncestor: '#jp_container_1',
     cssSelector: {
        videoPlay: '.jp-video-play',
        play: '.jp-play',
        pause: '.jp-pause',
        stop: '.jp-stop',
        seekBar: '.jp-seek-bar',
        playBar: '.jp-play-bar',
        mute: '.jp-mute',
        unmute: '.jp-unmute',
        volumeBar: '.jp-volume-bar',
        volumeBarValue: '.jp-volume-bar-value',
        volumeMax: '.jp-volume-max',
        playbackRateBar: '.jp-playback-rate-bar',
        playbackRateBarValue: '.jp-playback-rate-bar-value',
        currentTime: '.jp-current-time',
        duration: '.jp-duration',
        title: '.jp-title',
        fullScreen: '.jp-full-screen',
        restoreScreen: '.jp-restore-screen',
        repeat: '.jp-repeat',
        repeatOff: '.jp-repeat-off',
        gui: '.jp-gui',
        noSolution: '.jp-no-solution'
     },
     errorAlerts: true,
     warningAlerts: true
    });
})(jQuery);
