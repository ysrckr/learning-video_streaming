import './App.css';

import { Suspense, createResource } from 'solid-js';

const getVideoURL = async () => {
  try {
    const res = await fetch(
      'http://localhost:8000/videos?video_name=video.mp4',
    );
    return res.json();
  } catch (error) {
    console.error(error);
  }
};

function App() {
  const [url] = createResource(getVideoURL);

  console.log(url(), url.loading);

  return (
    <main>
      <Suspense fallback={<p>Loading</p>}>
        <video controls width={500}>
          <source src={url()?.['video_url']} type="video/mp4" />{' '}
          <p>Your browser cannot play the provided video file.</p>
        </video>
      </Suspense>
    </main>
  );
}

export default App;
