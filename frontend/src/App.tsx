import './App.css';

import { Suspense, createResource } from 'solid-js';

const getVideoURL = async () => {
  try {
    const res = await fetch('http://localhost:8000/videos?video_name=video');
    return res.json();
  } catch (error) {
    console.error(error);
  }
};

function App() {
  const [url] = createResource(getVideoURL);
  const videoURL = url()?.['video_url'];
  return (
    <main>
      <Suspense fallback={<p>Loading</p>}>
        <video controls loop width={500}>
          <source src={videoURL} type="video/mp4" />{' '}
          <p>Your browser cannot play the provided video file.</p>
        </video>
      </Suspense>
    </main>
  );
}

export default App;
