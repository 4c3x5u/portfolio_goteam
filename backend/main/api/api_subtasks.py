from rest_framework.views import APIView
from rest_framework.response import Response
from ..serializers.ser_subtask import SubtaskUpdateSerializer


class Subtasks(APIView):
    """
    Only used for updating subtasks.
    """
    @staticmethod
    def patch(request):
        serializer = SubtaskUpdateSerializer()
        validated_data = serializer.validate({
            'id': request.query_params.get('id'),
            'data': request.data,
            'user': {
                'username': request.META.get('HTTP_AUTH_USER'),
                'token': request.META.get('HTTP_AUTH_TOKEN')
            }
        })
        update_res = serializer.update(validated_data.get('instance'),
                                       validated_data.get('data'))
        return Response(update_res, 200)
