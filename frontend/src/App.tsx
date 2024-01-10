import './App.css';

function App() {
  return (
    <main>
      <video controls autoplay>
        <source src="http://localhost:8000/videos" type="video/mp4" />{' '}
        <p>Your browser cannot play the provided video file.</p>
      </video>
    </main>
  );
}

export default App;
