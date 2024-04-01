import React, { useState } from 'react';

function App() {
  const [text, setText] = useState('');
  const [response, setResponse] = useState('');
  const [error, setError] = useState(null);

  const checkText = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'text/plain'
        },
        body: text
      });

      if (!res.ok) {
        const errorData = await res.text();
        console.log(errorData);
        throw new Error(errorData); // Pass the error data to the Error constructor
      }

      const data = await res.json();
      console.log(data);
      setResponse(JSON.stringify(data)); // Convert the data to a string before setting it
      setError(null);
    } catch (error) {
      console.log(error);
      setError(error.toString()); // Convert the error to a string before setting it
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
      <h1>JSON Validator</h1>
      <textarea 
        value={text} 
        onChange={e => setText(e.target.value)} 
        style={{ width: '50%', height: '200px' }} 
      />
      <button onClick={checkText}>Sprawdź</button>
      <p>Response(AWS::IAM::Role Policy): {response}</p>
      {error && <p> {error}</p>}
    </div>
  );
}

export default App;