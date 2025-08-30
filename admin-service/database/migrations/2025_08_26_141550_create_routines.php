<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('tasks', function (Blueprint $table) {
            $table->id();
            $table->foreignId("habit_id")->references("id")->on("habits");
            $table->foreignId("user_id")->references("id")->on("users");
            $table->foreignId("packet_id")->references("id")->on("packets");
            $table->string("name");
            $table->text("description");
            $table->enum("difficulty", ["hard", "normal", "easy"]);
            $table->boolean("completed")->default(false);
            $table->timestamp('created_at')->default(DB::raw('CURRENT_TIMESTAMP'));
            $table->timestamp('updated_at')->nullable(true);
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('routines');
    }
};
