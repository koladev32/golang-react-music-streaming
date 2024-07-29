import React from 'react';
import AudioPlayer from 'react-h5-audio-player';
import 'react-h5-audio-player/lib/styles.css';

const MusicPlayer = ({ url }) => {
  return (
    <AudioPlayer
      autoPlay
      src={url}
      onPlay={e => console.log("Playing")}
      // other props here
    />
  );
};

export default MusicPlayer;
