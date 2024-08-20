import React, { useEffect, useState } from 'react';
import StateDisplay from './StateDisplay';

interface AgentState {
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
  FullPrompt: any[] | null;
  PostRationalisation: string;
  Error: string | null;
}

const GameSimulation: React.FC = () => {
  // Game state is a map of turn ID to a map of agent ID to agent turn
  const [gameState, setGameState] = useState<{ [turnID: number]: { [agentID: number]: AgentTurn } }>({});
  const [turn, setTurn] = useState(0);

  useEffect(() => {
    let url = new URL('ws://localhost:8080/ws');
    const ws = new WebSocket(url);

    ws.onopen = () => {
      console.log('WebSocket is open now.');
    };

    ws.onmessage = (evt) => {
      let turnData = JSON.parse(evt.data);

      setGameState((prevState) => {
        console.log('Received turn data:', turnData);
        return {
          ...prevState,
          [turn]: turnData
        };
      });

      setTurn(turn + 1);
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

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Game State</h1>
      <div className="grid grid-row-1 gap-4">
        {Object.entries(gameState).map(([turnID, turns]) => (
          <div key={turnID}>
            <h2 className="text-xl font-bold mb-2">Turn {turnID}</h2>
            <div className="grid grid-cols-1 gap-4">
              {Object.values(turns).map((agentTurn) => (
                <StateDisplay key={agentTurn.AgentID} agentTurn={agentTurn} />
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default GameSimulation;
