from rest_framework.decorators import api_view
from rest_framework.response import Response
from ..models import Column
from ..serializers.ser_task import TaskSerializer
from ..serializers.ser_subtask import SubtaskSerializer


@api_view(['POST'])
def tasks(request):
    column_tasks = Column.objects.filter(id=request.data.get('column'))
    task_serializer = TaskSerializer(
        data={'title': request.data.get('title'),
              'description': request.data.get('description'),
              'order': len(column_tasks) + 1,
              'column': request.data.get('column')}
    )
    if not task_serializer.is_valid():
        return Response(task_serializer.errors, 400)
    task = task_serializer.save()

    subtasks = request.data.get('subtasks')
    if subtasks:
        for i, title in enumerate(subtasks):
            subtask_serializer = SubtaskSerializer(
                data={'title': title.get('title'),
                      'order': i,
                      'task': task.id}
            )
            if not subtask_serializer.is_valid():
                task.delete()
                return Response(subtask_serializer.errors, 400)
            subtask_serializer.save()

    return Response({
        'msg': 'Task creation successful.',
        'task_id': task.id
    }, 201)

