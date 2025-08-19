<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('statistics', function (Blueprint $table) {
            $table->id();
            $table->foreignId("user_id")->references("id")->on("users");
            $table->integer("challenges")->default(0);
            $table->integer("events")->default(0);
            $table->integer("quests")->default(0);
            $table->integer("treasures")->default(0);
            $table->integer("tree_grown")->default(0);
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('statistics');
    }
};
