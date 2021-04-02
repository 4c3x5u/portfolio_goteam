from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Subtask
from ..serializers.ser_subtask import SubtaskSerializer


@api_view(['PATCH'])
def subtasks(request):
    subtask_id = request.data.get('id')
    if not subtask_id:
        return Response({
            'id': ErrorDetail(string='Subtask ID cannot be empty.',
                              code='blank')
        }, 400)

    data = request.data.get('data')
    if not data:
        return Response({
            'data': ErrorDetail(string='Data cannot be empty.', code='blank')
        }, 400)

    if 'title' in list(data.keys()) and not data.get('title'):
        return Response({
            'data.title': ErrorDetail(string='Title cannot be empty.',
                                      code='blank')
        }, 400)

    done = data.get('done')
    if 'done' in list(data.keys()) and (done == '' or done is None):
        return Response({
            'data.done': ErrorDetail(string='Done cannot be empty.',
                                     code='blank')
        }, 400)

    order = data.get('order')
    if 'order' in list(data.keys()) and (order == '' or order is None):
        return Response({
            'data.order': ErrorDetail(string='Order cannot be empty.',
                                      code='blank')
        }, 400)

    serializer = SubtaskSerializer(Subtask.objects.get(id=subtask_id),
                                   data=data,
                                   partial=True)
    if not serializer.is_valid():
        return Response(serializer.errors, 400)

    subtask = serializer.save()
    return Response({
        'msg': 'Subtask update successful.',
        'id': subtask.id
    }, 200)

