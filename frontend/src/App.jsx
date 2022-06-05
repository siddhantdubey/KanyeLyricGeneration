import styles from './App.module.css';
import { HopeProvider, Button, Heading } from '@hope-ui/solid'
import { createSignal } from 'solid-js';
import axios from 'axios';


function App() {
  const [lyrics, setLyrics] = createSignal([]);
  return (
    <HopeProvider>
      <div className={styles.App}>
        <div className={styles.main}>
          <div className={styles.top}>
            <header className={styles.AppHeader}>
              <Heading size="6xl" className={styles.heading}>Kanye Lyrics Generator</Heading>
            </header>
            <Button onClick={
              () => {
                axios.get("http://localhost:8080/sentence").then(
                  (response) => {
                    setLyrics(response.data.split("\n"))
                  }
                )
              }
            }>Generate Lyrics</Button>
          </div>
          <div className={styles.lyrics}>
            {/* iterate through lyrics */}
            {lyrics().map((line, index) => {
              return <p key={index}>{line}</p>
            })}
          </div>
        </div>
      </div>
    </HopeProvider>
  );
}

export default App;
