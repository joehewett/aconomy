# Aconomy: A Strategic Resource Management and Negotiation Game for LLMs

Aconomy is a turn-based strategy game that challenges AI agents to manage resources, negotiate with each other, and compete to be the first to accumulate a fixed number of gold. The game is designed to evaluate the abilities of language models in a multi-agent environment, focusing on resource management, strategic planning, and negotiation skills.

Things I'd like to see agents doing:
- Defrauding eachother by promising resources and not delivering
- Forming alliances to take down the leader
- Forming cartels to control the market
- Using the chat to manipulate other agents, e.g. by lying about their resources
- Overcoming global shortages by pooling resources to allow one agent to purchase a building and then share the output
- Using the chat to coordinate actions, e.g. "I'll buy a mine if you buy a farm"

### Future Features 
- **Constraints**: 
- **Skills**: Agents can have different skills that affect their abilities, e.g. negotiation, deception, resource management.
- **More Buildings**: Introduce new buildings with unique effects and resource requirements.
- **More Resources**: Add new resources with different uses and trade values.
- **More Actions**: Allow agents to perform additional actions, such as attacking other agents or sabotaging buildings.
- **Different Models**: Implement different AI models to compete against each other. Currently the model is fixed for all agents. 

## Game Mechanics

### Resources
- **Gold**: The primary victory resource.
- **Wheat**: Used to feed workers and can be traded.

### Buildings
- **Farm**: Produces a number of wheat per turn when manned.
- **Mine**: Produces a number of gold per turn when manned.

### Workers
- Cost gold to recruit.
- Consume wheat per turn.
- Required to operate buildings.
- Flee if not fed at the start of a turn.

### Key Rules
- Wheat decays by a percentage per turn.
- Each agent has a fixed number of actions per turn.
- The game ends when an agent reaches 1000 gold or after 100 turns.

### Actions
1. Give resources (gold or wheat) to another agent.
2. Buy workers
3. Buy buildings

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/Aconomy.git
   ```

2. Navigate to the project directory:
   ```
   cd Aconomy
   ```

3. Build the project:
   ```
   go build
   ```

4. Run the game:
   ```
   ./Aconomy
   ```

## Contributing

Contributions to Aconomy are welcome! Please feel free to submit pull requests, create issues, or suggest enhancements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
