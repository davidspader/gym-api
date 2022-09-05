<?php

namespace App\Http\Controllers\Auth\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;

class LoginController extends Controller
{
    public function login(Request $request)
    {
        $credentials = $request->only('email', 'password');

        if(!auth()->attempt($credentials)) {
            return response()->json( [
                'error' => [
                    'message' => 'Invalid credentials'
                ]
            ], 401);
        }

        $token = auth()->user()->createToken('auth_token');

        return response()->json([
            'data' => [
                'user' => auth()->user(),
                'token' => $token->plainTextToken
            ]
        ]);
    }

    public function logout()
    {
        auth()->user()->currentAccessToken()->delete();

        return response()->json([], 204);
    }
}
