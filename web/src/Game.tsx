import React, { useEffect, useState } from 'react';
import TurnDisplay from './TurnDisplay';
import OpenAI from 'openai';

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

  useEffect(() => {
    let url = new URL('ws://localhost:8080/ws');
    const ws = new WebSocket(url);

    ws.onopen = () => {
      console.log('WebSocket is open now.');
    };

    ws.onmessage = (evt) => {
      let turnData: AgentTurn = JSON.parse(evt.data);
      console.log('Received:', turnData);
      setGameTurns(prevTurns => [...prevTurns, turnData]);
    };

    ws.onerror = (err) => {
      console.log('WebSocket error:', err);
    };

    ws.onclose = () => {
      // Submit a close message to the server
      ws.close();
      console.log('WebSocket is closed now.');
    };

    return () => {
      ws.close();
    };
  }, []);

  if (gameState.length === 0) {
    return (
      <div className="container mx-auto p-4">
        <h1 className="text-2xl font-bold">Game State</h1>
        <p className="text-lg">Waiting for game state...</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Game State</h1>
      <div className="grid grid-row-1 gap-4">
        {Object.entries(gameState).map(([turnID, turns]) => (
          <TurnDisplay key={turnID} agentTurn={turns} />
        ))}
      </div>
    </div>
  );
};

export default GameSimulation;
