import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';
import React, { useContext } from 'react';
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

  const handleOnDragEnd = async (result) => {
    if (!result?.source?.droppableId || !result?.destination?.droppableId) {
      return;
    }

    // parse source and destination droppable IDs so that they are comparable
    // to the column order
    const iSource = parseInt(result.source.droppableId)
    const iDest = parseInt(result.destination.droppableId)

    // find the source column – the one that the task is initially in
    const source = activeBoard.columns[iSource]

    // pop the task that's being moved out of the source
    const [item] = source.tasks.splice(result.source.index, 1);

    // update source tasks' orders
    const sourceTasks = source.tasks.map((task, index) => ({
      ...task, colNo: iSource, order: index,
    }));

    // find the destination column – the that the task is being moved into
    const destination = activeBoard.columns[iDest]

    // insert the task that's being moved into the destination
    destination.tasks.splice(result.destination.index, 0, item);

    // update destination tasks' orders
    const destinationTasks = destination.tasks.map((task, index) => ({
      ...task, colNo: iDest, order: index,
    }));

    // update client state ahead of API calls for faster UI - if errors occur
    // during API calls, an error will be displayed and the state will be reset
    setActiveBoard({
      ...activeBoard,
      columns: activeBoard.columns.map((column, i) => {
        switch (i) {
          case iDest:
            return { tasks: destinationTasks };
          case iSource:
            return { tasks: sourceTasks };
          default:
            return column
        }
      }),
    });

    try {
      // update the source and the destination in the database
      let sourceProm = sourceTasks.length > 0 && iSource !== iDest
        ? TasksAPI.patch(sourceTasks)
        : Promise.resolve()
      let destProm = destinationTasks.length > 0
        ? await TasksAPI.patch(destinationTasks)
        : Promise.resolve()
      await Promise.all([sourceProm, destProm])
    } catch (err) {
      notify(
        'Unable to update tasks.',
        `${err?.message || 'Server Error'}.`,
      );
      setIsLoading(true);
    } finally {
      loadBoard(activeBoard.id);

    }
  };

  return (
    <Row id="Board">
      <DragDropContext onDragEnd={handleOnDragEnd}>
        {activeBoard.columns
          && activeBoard.columns.map((column, i) => (
            <Column
              id={i}
              key={i}
              order={i}
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
