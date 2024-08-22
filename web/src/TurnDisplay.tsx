import React from 'react';
import { AgentTurn } from './Game';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./components/ui/card"
import { Label } from "./components/ui/label"
import { PromptModal } from "./components/PromptModal"

interface TurnDisplayProps {
  agentTurn: AgentTurn;
}

const TurnDisplay: React.FC<TurnDisplayProps> = ({ agentTurn }) => {
  function emojiString(count: number, emoji: string) {
    return Array(count).fill(emoji).join('');
  }

  function buildingsString(buildings: { Type: string, Manned: boolean }[]) {
    return buildings.map(building => {
      return `${building.Type === 'Farm' ? '🚜' : '⛏️'} ${building.Manned ? '(manned)' : '(unmanned)'}`;
    }).join(', ');
  }

  return (
    <Card className="col-span-1 font-mono text-left">
      <CardHeader>
        <CardTitle>Agent {agentTurn.AgentID} / / Turn {agentTurn.Turn}</CardTitle>
        <CardDescription className="space-y-4">
          <span className="text-lg">
            Action: {agentTurn.Action}
          </span>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid w-full items-start gap-4">
          <div className="flex flex-col space-y-1.5 items-start text-left">
            <p>Gold: 🪙 {agentTurn.StartState.Gold} --&gt; {agentTurn.EndState.Gold}</p>
            <p>Wheat: {emojiString(agentTurn.EndState.Wheat, '🌾')} {agentTurn.StartState.Wheat} --&gt; {agentTurn.EndState.Wheat}</p>
            <p>Workers: {emojiString(agentTurn.StartState.Workers, '👷')} {agentTurn.StartState.Workers} --&gt; {agentTurn.EndState.Workers}</p>
            {(agentTurn.StartState.Buildings && agentTurn.EndState.Buildings) && (
              <p>Buildings: {buildingsString(agentTurn.EndState.Buildings)} {agentTurn.StartState.Buildings.length} --&gt; {agentTurn.EndState.Buildings.length}</p>
            )}
          </div>
          <div className="flex flex-col space-y-1.5">
            <Label>🧠 Strategy</Label>
            <p>{agentTurn.Strategy}</p>
          </div>
          <div className="flex flex-col space-y-1.5">
            <Label>💭 Post Rationalisation</Label>
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
        {agentTurn.FullPrompt && (
          <PromptModal agentTurn={agentTurn} />
        )}

      </CardFooter>
    </Card>
  );
};

export default TurnDisplay;

