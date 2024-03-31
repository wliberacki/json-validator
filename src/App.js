import React, { useState } from 'react';

function App() {
  const [text, setText] = useState('');
  const [response, setResponse] = useState('');

  const checkText = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ text })
      });
  
      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }
  
      const data = await res.json();
      console.log(data.response);
      setResponse(data.response);
    } catch (error) {
      console.error('An error occurred:', error);
    }
  };

  return (
    <div>
      <input type="text" value={text} onChange={e => setText(e.target.value)} />
      <button onClick={checkText}>Sprawdź</button>
      <p>{response}</p>
    </div>
  );
}

export default App;