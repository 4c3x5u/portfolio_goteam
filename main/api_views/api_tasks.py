from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Column
from ..serializers.ser_task import TaskSerializer
from ..serializers.ser_subtask import SubtaskSerializer


@api_view(['POST'])
def tasks(request):
    column_id = request.data.get('column')
    if not column_id:
        return Response({
            'column': ErrorDetail(string='Column cannot be empty.',
                                  code='blank')
        }, 400)

    task_serializer = TaskSerializer(
        data={'title': request.data.get('title'),
              'description': request.data.get('description'),
              'order': Column.objects.filter(id=column_id).count() + 1,
              'column': request.data.get('column')}
    )
    if not task_serializer.is_valid():
        return Response(task_serializer.errors, 400)
    task = task_serializer.save()

    subtasks = request.data.get('subtasks')
    if subtasks:
        for i, subtask in enumerate(subtasks):
            subtask_serializer = SubtaskSerializer(
                data={'title': subtask.get('title'),
                      'order': i,
                      'task': task.id}
            )
            if not subtask_serializer.is_valid():
                task.delete()
                return Response({'subtask': subtask_serializer.errors}, 400)
            subtask_serializer.save()

    return Response({
        'msg': 'Task creation successful.',
        'task_id': task.id
    }, 201)

