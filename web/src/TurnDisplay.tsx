// src/StateDisplay.tsx
import React, { useState } from 'react';
import { AgentTurn } from './Game';
import { Button } from "./components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./components/ui/card"
import { Input } from "./components/ui/input"
import { Label } from "./components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./components/ui/select"
import { PromptModal } from "./components/PromptModal"



interface TurnDisplayProps {
  agentTurn: AgentTurn;
}


const TurnDisplay: React.FC<TurnDisplayProps> = ({ agentTurn }) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const toggleExpanded = () => {
    setIsExpanded(!isExpanded);
  };

  function emojiString(count: number, emoji: string) {
    return Array(count).fill(emoji).join('');
  }

  return (
    <Card className="w-[700px] font-mono text-left">
      <CardHeader>
        <CardTitle>Agent {agentTurn.AgentID} // Turn {agentTurn.Turn}</CardTitle>
        <CardDescription>Action: {agentTurn.Action}</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid cols-2 w-full items-start gap-4">
          <div className="flex flex-col space-y-1.5 items-start text-left">
            <p>Gold: 🪙 {agentTurn.StartState.Gold}</p>
            <p>Wheat: {emojiString(agentTurn.StartState.Wheat, '🌾')}</p>
            <p>Workers: {emojiString(agentTurn.StartState.Workers, '👷')}</p>
            {agentTurn.StartState.Buildings.length > 0 && (
              <Label>🏠 Buildings:</Label>
            )}
            {agentTurn.StartState.Buildings.map((building, idx) => (
              <p key={idx}>
                {building.Type} {building.Manned ? '👨‍🌾' : '❌'}
              </p>
            ))}
          </div>
          <div className="flex flex-col space-y-1.5">
            <p>{agentTurn.Strategy}</p>
          </div>
          <div className="flex flex-col space-y-1.5 items-start text-left">
            <p>Gold: 🪙 {agentTurn.EndState.Gold}</p>
            <p>Wheat: {emojiString(agentTurn.EndState.Wheat, '🌾')}</p>
            <p>Workers: {emojiString(agentTurn.EndState.Workers, '👷')}</p>
            <Label>🏠 Buildings:</Label>
            {agentTurn.EndState.Buildings.map((building, idx) => (
              <p key={idx}>
                {building.Type} {building.Manned ? '👨‍🌾' : '❌'}
              </p>
            ))}
          </div>
          <div className="flex flex-col space-y-1.5">
            <Label>💭 Post Rationalisation:</Label>
            <p>{agentTurn.PostRationalisation}</p>
          </div>
          {agentTurn.Error && (
            <div className="flex flex-col space-y-1.5">
              <Label>⚠️ Error:</Label>
              <p className="text-red-500">{agentTurn.Error}</p>
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button variant="outline">Close</Button>

        <PromptModal agentTurn={agentTurn} />

      </CardFooter>
    </Card>
  );
};

export default TurnDisplay;

