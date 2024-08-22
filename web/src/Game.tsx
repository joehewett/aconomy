import React, { useEffect, useState } from 'react';
import TurnDisplay from './TurnDisplay';
import OpenAI from 'openai';
import { Input } from './components/ui/input';
import { ActionBar } from './components/ActionBar';

export interface AgentState {
  Gold: number;
  Wheat: number;
  Workers: number;
  Buildings: { Type: string, Manned: boolean }[];
}

export interface AgentTurn {
  AgentID: number;
  StartState: AgentState;
  Strategy: string;
  Action: string;
  EndState: AgentState;
  FullPrompt: OpenAI.ChatCompletionMessage[];
  PostRationalisation: string;
  Error: string | null;
  Turn: number;
}

const GameSimulation: React.FC = () => {
  const [gameState, setGameTurns] = useState<AgentTurn[]>([]);
  const [started, setStarted] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [apiKey, setAPIKey] = useState<string>('');
  const [websocket, setWebsocket] = useState<WebSocket | null>(null);

  function storeAPIKey(apiKey: string) {
    setAPIKey(apiKey);
    localStorage.setItem('openai_api_key', apiKey);
  }

  useEffect(() => {
    let storedAPIKey = localStorage.getItem('openai_api_key');
    if (storedAPIKey) {
      setAPIKey(storedAPIKey);
    }
  }, []);

  function StartGame() {
    setGameTurns([]);
    setLoading(true);
    setError(null)

    console.log('env var: ', process.env.RAILWAY_PUBLIC_DOMAIN);

    let host = process.env.SERVER_HOST;

    if (!host) {
      setError('Server host not set. Please check the environment variables.');
      setLoading(false);
      return;
    }

    let url = new URL(host);
    url.searchParams.append('api_key', apiKey);

    const ws = new WebSocket(url);

    setWebsocket(ws);

    ws.onopen = () => {
      console.log('WebSocket is open now.');
    };

    ws.onmessage = (evt) => {
      setStarted(true);
      setLoading(false);
      let turnData: AgentTurn = JSON.parse(evt.data);
      console.log('Received:', turnData);
      setGameTurns(prevTurns => [...prevTurns, turnData]);
    };

    ws.onerror = (err) => {
      console.log('WebSocket error:', err);
      setLoading(false);
      setError(`Failed to connect to the server. The server might be offline. Please try again or come back later.`);
    };

    ws.onclose = () => {
      // Submit a close message to the server
      ws.close();
      console.log('WebSocket is closed now.');
      setStarted(false);
    };

    return () => {
      ws.close();
    };
  }

  function stopGame() {
    setStarted(false);
    if (websocket) {
      websocket.close();
    }
  }

  return (
    <div className="container mx-auto py-32 font-mono">
      <div className="flex flex-row items-center gap-4 justify-between">
        <div>
          <h2 className="text-left  font-serif py-4 leading-7 text-gray-900 sm:truncate sm:text-7xl sm:tracking-tight">aconomy🌾</h2>
          {/* Subtitle */}
          <p className="text-gray-600 text-left"></p>

        </div>
        <Input type="text" className="w-64" placeholder="OpenAI API Key" onChange={(e) => storeAPIKey(e.target.value)} value={apiKey} />
        <ActionBar StartGame={StartGame} loading={loading} stopGame={stopGame} started={started} />
      </div>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mt-4" role="alert">
          <strong className="font-bold">Error: </strong>
          <span className="block sm:inline">{error}</span>
        </div>
      )}

      <div className="container px-0 py-4">
        <div className="grid grid-cols-3 gap-4">
          {Object.entries(gameState).map(([turnID, turns]) => (
            <TurnDisplay key={turnID} agentTurn={turns} />
          ))}

        </div>
      </div>
    </div>
  );
};

export default GameSimulation;
