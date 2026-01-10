import { Greet } from '../wailsjs/go/main/App';

function App() {
  const action = async () => {
    const res = await Greet('Adam');
    alert(res);
  };

  return (
    <div className="min-h-screen bg-white grid grid-cols-1 place-items-center justify-items-center mx-auto py-8">
      <div className="text-blue-900 text-2xl font-bold font-mono">
        <h1 className="content-center">Vite + React + TS + Tailwind</h1>
      </div>
      <button onClick={action}>Action</button>
    </div>
  );
}

export default App;
