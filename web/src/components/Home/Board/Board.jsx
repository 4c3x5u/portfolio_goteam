import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';
import { DragDropContext } from 'react-beautiful-dnd';

import AppContext from '../../../AppContext';
import TasksAPI from '../../../api/TasksAPI';
import Column from './Column/Column';

import './board.sass';

const Board = ({ handleActivate }) => {
  const {
    activeBoard,
    loadBoard,
    setIsLoading,
    notify,
    setActiveBoard,
  } = useContext(AppContext);

  console.log("~~~ BOARD CALL")
  console.log("ACTIVE BOARD COLUMNS: " + JSON.stringify(activeBoard.columns))

  const handleOnDragEnd = async (result) => {
    if (!result?.source?.droppableId || !result?.destination?.droppableId) {
      return;
    }
    try {
      // Find the "source" column – the one that the task is initially in
      const source = activeBoard.columns.find((column) => (
        column.id.toString() === result.source.droppableId
      ));

      // Pop the tasks that's being moved out of the "source"
      const [item] = source.tasks.splice(result.source.index, 1);

      // Update "source" tasks' orders
      const sourceTasks = source.tasks.map((task, index) => ({
        ...task,
        column: source.id,
        order: index,
      }));

      // Find the "destination" column – the that the task is being moved into
      const destination = activeBoard.columns.find((column) => (
        column.id.toString() === result.destination.droppableId
      ));

      // Insert the task that's being moved into the "destination"
      destination.tasks.splice(result.destination.index, 0, item);

      // Update "destination" tasks' orders
      const destinationTasks = destination.tasks.map((task, index) => ({
        ...task,
        column: destination.id,
        order: index,
      }));

      // Update client state (to avoid server-response wait time)
      setActiveBoard({
        ...activeBoard,
        columns: activeBoard.columns.map((column) => {
          if (column.id === destination.id) {
            return { ...column, tasks: destinationTasks };
          }
          if (column.id === source.id) {
            return { ...column, tasks: sourceTasks };
          }
          return column;
        }),
      });

      // Update the "source" and the "destination" in the database
      await TasksAPI.patch(sourceTasks);
      await TasksAPI.patch(destinationTasks);
    } catch (err) {
      notify(
        'Unable to update task.',
        `${err?.message || 'Server Error'}.`,
      );
      setIsLoading(true);
      loadBoard();
    }
  };

  return (
    <Row id="Board">
      <DragDropContext onDragEnd={handleOnDragEnd}>
        {activeBoard.columns
          && activeBoard.columns.map((column, i) => (
            <Column
              key={i}
              order={i + 1}
              tasks={column.tasks}
              handleActivate={handleActivate}
            />
          ))}
      </DragDropContext>
    </Row>
  );
};

Board.propTypes = { handleActivate: PropTypes.func.isRequired };

export default Board;
