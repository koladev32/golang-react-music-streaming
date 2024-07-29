"use client"

import React, { useEffect, useState } from 'react';
import MusicPlayer from "@/components/music-player";

export default function Home() {
  const [songs, setSongs] = useState([]);
  const [currentSong, setCurrentSong] = useState(null);

  useEffect(() => {
    fetch('/api/songs')
      .then((response) => response.json())
      .then((data) => setSongs(data));
  }, []);

  const playSong = (song) => {
    setCurrentSong(song.url);
  };

  return (
    <div className="container mx-auto px-4">
      <h1 className="text-4xl font-bold text-center my-8">Music Streaming</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {songs.map((song) => (
          <div key={song.id} className="bg-white p-4 rounded-lg shadow-md">
            <img src={song.thumbnail} alt={song.name} className="w-full h-48 object-cover rounded-lg mb-4" />
            <h2 className="text-2xl font-semibold mb-2">{song.name}</h2>
            <p className="text-gray-600 mb-2">By {song.author}</p>
            <button onClick={() => playSong(song)} className="text-blue-500 hover:underline">
              Listen Now
            </button>
          </div>
        ))}
      </div>
      {currentSong && (
        <div className="fixed bottom-0 left-0 right-0 bg-gray-800 p-4">
          <MusicPlayer url={currentSong} />
        </div>
      )}
    </div>
  );
}
