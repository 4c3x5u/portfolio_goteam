from rest_framework.generics import CreateAPIView
from rest_framework.exceptions import ValidationError
from django.core.exceptions import ObjectDoesNotExist
from .serializers import UserSerializer
from .models import Team


class UserCreate(CreateAPIView):
    serializer_class = UserSerializer

    def create(self, request, *args, **kwargs):
        pw = request.data.get('password')
        cf = request.data.get('password_confirmation')
        if pw == cf:
            serializer = UserSerializer(data=request.data, many=False)
            serializer.is_valid()
            try:
                Team.objects.get(id=request.data.get('team'))
            except ObjectDoesNotExist:
                team = Team.objects.create()
                request.data.update({'team': team.id})
            return super().create(request, *args, **kwargs)
        return ValidationError({
            'password_confirmation': 'Password confirmation must match the'
                                     'password'
        })


# @api_view(['POST'])
# def team(request):
#     return HttpResponse(Team.objects.create())
