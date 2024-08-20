// src/StateDisplay.tsx
import React from 'react';
import { AgentTurn } from './Game';

interface StateDisplayProps {
  agentTurn: AgentTurn;
}

const StateDisplay: React.FC<StateDisplayProps> = ({ agentTurn }) => {
  const { AgentID, StartState, EndState } = agentTurn;

  return (
    <div className="bg-white shadow-md rounded p-4">
      <h2 className="text-xl font-bold mb-2">Agent {AgentID}</h2>
      <div className="mb-2">
        <strong>Start State:</strong>
        <pre className="bg-gray-100 p-2 rounded">{JSON.stringify(StartState, null, 2)}</pre>
      </div>
      <div className="mb-2">
        <strong>End State:</strong>
        <pre className="bg-gray-100 p-2 rounded">{JSON.stringify(EndState, null, 2)}</pre>
      </div>
    </div>
  );
};

export default StateDisplay;
